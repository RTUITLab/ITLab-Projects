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
	var issues []milestone.IssuesWithMilestoneID
	values := r.URL.Query()

	start := getUint(values, "start")

	count := getUint(values, "count")
	if count == 0 {
		count = uint64(a.Repository.Repo.Count())
	}

	ctx := context.Background()

	if err := a.Repository.Issue.GetAllFiltered(
		ctx,
		bson.M{
			"state": "open",
		},
		func(c *mongo.Cursor) error {
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