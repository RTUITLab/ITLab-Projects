package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	e "github.com/ITLab-Projects/pkg/err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/models/tag"
)

// GetTags
//
// @Summary return all tags
//
// @Produce json
//
// @Tags tags
//
// @Description return all tags
//
// @Router /api/v1/projects/tags [get]
//
// @Success 200 {array} tag.Tag
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
func (a *Api) GetTags(w http.ResponseWriter, r *http.Request) {
	var tags []tag.Tag

	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()
	
	if err := a.Repository.Tag.GetAllFiltered(
		ctx,
		bson.M{},
		func(c *mongo.Cursor) error {
			if err := c.All(
				ctx,
				&tags,
			); err != nil {
				return err
			}

			return c.Err()
		},
	); err == mongo.ErrNoDocuments {
		// pass
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Faield to get tags",
			},
		)
		prepare("GetTags", err).Error()
		return
	}

	json.NewEncoder(w).Encode(tags)
}