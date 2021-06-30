package issues

import (
	"strings"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/go-kit/kit/log/level"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
	"github.com/go-kit/kit/log"
	e "github.com/ITLab-Projects/pkg/err"
	"errors"
)


func init() {
	// to generate swagger
	_ = e.Message{}
}

var (
	ErrFailedToGetIssues 	= errors.New("Failed to get issues")
	ErrFailedToGetLabels	= errors.New("Failed to get issues labels")
)

type ServiceImp struct {
	repository 	Repository
	logger		log.Logger
}

func New(
	repository 	Repository,
	logger 		log.Logger,
) *ServiceImp {
	return &ServiceImp{
		repository: repository,
		logger: logger,
	}
}

// GetIssues
//
// @Summary return issues
//
// @Tags issues
//
// @Security ApiKeyAuth
// 
// @Produce json
//
// @Description return issues according to query params
//
// @Router /v2/issues [get]
//
// @Param start query integer false "represent how mush skip first issues"
//
// @Param count query integer false "set limit of getting issues standart and max 50"
//
// @Param name query string false "search to name of issues, title of milestones and repository names"
//
// @Param tag query string false "search of label name of issues"
//
// @Success 200 {object} issues.GetIssuesResp
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
func (s *ServiceImp) GetIssues(
	ctx 		context.Context, 
	Query		GetIssuesQuery,
) ([]*milestone.IssuesWithMilestoneID, error) {
	logger := log.With(s.logger, "method", "GetIssues")
	filter, err := s.buildFilterForGetIssues(
		ctx,
		Query.Name,
		Query.Tag,
	)
	if err != nil {
		level.Error(logger).Log("Failed to get issues: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetIssues,
			http.StatusInternalServerError,
		)
	}
	if Query.Count == 0 || Query.Count > 50 {
		Query.Count = 50
	}

	is, err := s.repository.GetChunckedIssues(
		ctx,
		filter,
		bson.D{ {"createdat", -1}, {"deleted", 1}},
		int64(Query.Start),
		int64(Query.Count),
	)
	if err == mongo.ErrNoDocuments {
		is = []*milestone.IssuesWithMilestoneID{}
	} else if err != nil {
		level.Error(logger).Log("Failed to get issues: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetIssues,
			http.StatusInternalServerError,
		)
	}

	return is, nil
}

func (s *ServiceImp) buildFilterForGetIssues(
	ctx 		context.Context,
	name, tag 	string,
) (interface{}, error) {
	filter := bson.M{
		"state": "open",
	}

	if tag != "" {
		s.buildFilterByLabelTags(
			ctx,
			&filter,
			tag,
		)
	}

	if name != "" {
		if err := s.buildFilterByName(
			ctx,
			&filter,
			name,
		); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (s *ServiceImp) buildFilterByLabelTags(
	ctx 	context.Context,
	filter 	*bson.M,
	tag		string,
) {
	tags := strings.Split(tag, " ")

	(map[string]interface{})(*filter)["labels.name"] = bson.M{"$in": tags}
}

func(s *ServiceImp) buildFilterByName(
	ctx		context.Context,
	filter	*bson.M,
	name	string,
) error {
	type IDs struct {
		ID 	uint64		`bson:"id"`
	}

	var reposID	[]*IDs
	if err := s.repository.GetReposAndScanTo(
		ctx,
		bson.M{
			"name": bson.M{"$regex": name, "$options": "-i"},
		},
		&reposID,
		options.Find(),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	var milestonesID []IDs
	if err := s.repository.GetMilestonesAndScanTo(
		ctx,
		bson.M{
			"title": bson.M{"$regex": name, "$options": "-i"},
		},
		&milestonesID,
		options.Find(),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
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

	if len(reposID) > 0 {
		var ids []uint64
		for _, id := range reposID {
			ids = append(ids, id.ID)
		}

		array = append(
			array,
			bson.D {
				{"repo_id", bson.M{"$in": ids}},
			},
		)
	}

	if len(milestonesID) > 0 {
		var ids []uint64
		for _, id := range milestonesID {
			ids = append(ids, id.ID)
		}
		array = append(
			array, 
			bson.D {
				{"milestone_id", bson.M{"$in": ids}},
			},
		)
	}

	(map[string]interface{})(*filter)["$or"] = array

	return nil
}