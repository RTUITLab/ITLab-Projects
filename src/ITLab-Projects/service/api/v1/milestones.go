package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetIssues
//
// @Summary return issues
//
// @Tags issues
//
// @Produce json
//
// @Description return issues according to query params
//
// @Router /api/projects/issues [get]
//
// @Param start query integer false "represent how mush skip first issues"
//
// @Param count query integer false "set limit of getting issues"
//
// @Param name query string false "search to name of issues, title of milestones and repository names"
//
// @Param tag query string false "search of label name of issues"
//
// @Success 200 {array} milestone.IssuesWithMilestoneID
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
func (a *Api) GetIssues(w http.ResponseWriter, r *http.Request) {
	var issues []milestone.IssuesWithMilestoneID
	values := r.URL.Query()

	start := getUint(values, "start")

	count := getUint(values, "count")
	if count == 0 {
		count = uint64(a.Repository.Issue.Count())
	}

	filter := bson.M{
		"state": "open",
	}

	name := values.Get("name")
	tag := values.Get("tag")

	ctx, cancel := context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
	defer cancel()
	if tag != "" {
		filter = a.buildFilterForLabelTags(
			ctx,
			filter,
			tag,
		)
	}

	if name != "" {
		if _filter, err := a.buildFilterByNameForIssues(
			ctx,
			filter,
			name,
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				e.Message{
					Message: "Failed to find issues",
				},
			)
			prepare("GetIssues", err).Error()
			return
		} else {
			filter = _filter
		}

		logrus.Debug(filter)
	}

	if err := a.Repository.Issue.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			defer c.Close(ctx)
			if err := c.All(
				ctx,
				&issues,
			); err != nil {
				return err
			}

			return c.Err()
		},
		options.Find().
		SetSort(bson.D{ {"createdat", -1}, {"deleted", 1}} ).
		SetSkip(int64(start)).
		SetLimit(int64(count)),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Faile to get Issues",
			},
		)
		prepare("GetIssues", err).Error()
		return
	}

	json.NewEncoder(w).Encode(
		issues,
	)
}

func (a *Api) buildFilterByNameForIssues(ctx context.Context, filter bson.M, name string) (bson.M, error) {
	type IDs struct {
		ID 	uint64		`bson:"id"`
	}

	type ids []uint64
	

	var repoIDs []IDs
	if err := a.Repository.Repo.GetAllFiltered(
		ctx,
		bson.M{
			"name": bson.M{"$regex": name, "$options": "-i"},
		},
		func(c *mongo.Cursor) error {
			return c.All(
				ctx,
				&repoIDs,
			)
		},
		options.Find(),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}

	var milestoneIDs []IDs
	if err := a.Repository.Milestone.GetAllFiltered(
		ctx,
		bson.M{
			"title": bson.M{"$regex": name, "$options": "-i"},
		},
		
		func(c *mongo.Cursor) error {
			defer c.Close(ctx)
			return c.All(
				ctx,
				&milestoneIDs,
			)
		},
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}

	var array bson.A
	array = append(
		array,
		bson.D {
			{"title", bson.M{"$regex": name, "$options": "-i"}},
		},
		bson.D{
			{"description", bson.M{"$regex": name, "$options": "-i"}},
		},
	)

	if len(repoIDs) > 0 {
		var ids []uint64
		for _, id := range repoIDs {
			ids = append(ids, id.ID)
		}

		array = append(
			array,
			bson.D {
				{"repo_id", bson.M{"$in": ids}},
			},
		)
	}

	if len(milestoneIDs) > 0 {
		var ids []uint64
		for _, id := range milestoneIDs {
			ids = append(ids, id.ID)
		}
		array = append(
			array, 
			bson.D {
				{"milestone_id", bson.M{"$in": ids}},
			},
		)
	}

	f := func (m map[string]interface{}) bson.M {
		m["$or"] = array
		return m
	}(filter)

	return f, nil
}

func (a *Api) buildFilterForLabelTags(
	ctx context.Context, 
	filter bson.M, 
	tag string,
) (bson.M) {
	tags := strings.Split(tag, " ")

	return func(m map[string]interface{}) bson.M {
		m["labels.name"] = bson.M{"$in": tags}
		return m
	}(filter)
}

// GetLabels
//
// @Summary return labels
//
// @Tags issues
//
// @Produce json
//
// @Description return all unique labels of issues
//
// @Router /api/projects/issues/labels [get]
//
// @Success 200 {array} string
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
func (a *Api) GetLabels(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		1*time.Second,
	)
	defer cancel()

	names, err := a.Repository.Issue.Distinct(
		ctx,
		"labels.name",
		bson.M{"state": "open"},
		options.Distinct(),
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to get labels",
			},
		)
		prepare("GetLabels", err).Error()
		return
	}

	json.NewEncoder(w).Encode(names)
}