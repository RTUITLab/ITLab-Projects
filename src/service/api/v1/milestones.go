package v1

import (
	e "github.com/ITLab-Projects/pkg/err"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *Api) GetIssues(w http.ResponseWriter, r *http.Request) {
	var ms []milestone.MilestoneInRepo

	values := r.URL.Query()

	start := getUint(values, "start")

	count := getUint(values, "count")
	if count == 0 {
		count = uint64(a.Repository.Repo.Count())
	}

	ctx := context.Background()

	if err := a.Repository.Milestone.GetAllFiltered(
		ctx,
		bson.M{
			"state": "open",
		},
		func(c *mongo.Cursor) error {
			if err := c.All(
				ctx,
				&ms,
			); err != nil {
				return err
			}

			return c.Err()
		},
		options.Find().
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

	var is []milestone.IssuesWithMilestoneID

	if len(ms) == 0 {
		json.NewEncoder(w).Encode(
			is,
		)
		return
	}

	for _, m := range ms {
		for _, i := range m.Issues {
			if i.State == "closed" {
				continue
			}

			is = append(
				is, 
				milestone.IssuesWithMilestoneID{
					MilestoneID: m.ID,
					Issue: i,
				},
			)
		}
	}

	json.NewEncoder(w).Encode(
		is,
	)	
}