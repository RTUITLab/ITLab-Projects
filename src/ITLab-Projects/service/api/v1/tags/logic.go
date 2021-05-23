package tags

import (
	"context"
	"errors"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/statuscode"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/kit/log"
)

var (
	ErrGetTags		= errors.New("Faield to get tags")
)

type service struct {
	repository 	Repository
	logger		log.Logger
}

func New(
	Repository 	Repository,
	logger		log.Logger,
) *service {
	return &service{
		repository: Repository,
		logger: logger,
	}
}

// GetAllTags
//
// @Summary return all tags
//
// @Produce json
//
// @Tags tags
//
// @Description return all tags
//
// @Router /api/projects/tags [get]
//
// @Success 200 {array} tag.Tag
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
func (s *service) GetAllTags(
	ctx context.Context,
) ([]*tag.Tag, error) {
	logger := log.With(s.logger, "method", "GetAllTags")
	tgs, err := s.repository.GetAllTags(
		ctx,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		logger.Log("Failed to get tags: err", err)
		return nil, statuscode.WrapStatusError(
			ErrGetTags,
			http.StatusInternalServerError,
		)
	}

	return tgs, nil
}