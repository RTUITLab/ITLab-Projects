package tags

import (
	"context"

	"github.com/ITLab-Projects/pkg/models/tag"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/kit/log"
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

func (s *service) GetAllTags(
	ctx context.Context,
) ([]*tag.Tag, error) {
	tgs, err := s.repository.GetAllTags(
		ctx,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}

	return tgs, nil
}