package landing

import (
	"context"
	"errors"
	"net/http"
	"strings"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrTagNotFound				= 	errors.New("Tag not found")
	ErrFailedToGetLandings		= 	errors.New("Failed to get landings")
	ErrLandingNotFound			=	errors.New("Landing not found")
	ErrFaieldToGetLanding		= 	errors.New("Failed to get landing")
)

func init() {
	// to generate swagger
	_ = e.Message{}
}

type service struct {
	repository 	Repository
	logger		log.Logger
}

func New(
	repository 	Repository,
	logger		log.Logger,
) *service {
	return &service{
		repository: repository,
		logger: logger,
	}
}

func (s *service) GetAllLandings(
	ctx context.Context, 
	start int64, 
	count int64, 
	tag string, 
	name string,
) ([]*landing.LandingCompact, error) {
	logger := log.With(s.logger, "method", "GetAllLandings")
	filter, err := s.buildFilterForGetLanding(
		ctx,
		name,
		tag,
	)
	if err == ErrTagNotFound {
		return []*landing.LandingCompact{}, nil
	} else if err != nil {
		level.Error(logger).Log("Err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetLandings,
			http.StatusInternalServerError,
		)
	}

	ls, err := s.repository.GetFiltrSortLandingCompactFromTo(
		ctx,
		filter,
		bson.D{},
		start,
		count,
	)
	if err == mongo.ErrNoDocuments {
		return []*landing.LandingCompact{}, nil
	} else if err != nil {
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetLandings,
			http.StatusInternalServerError,
		)
	}

	return ls, nil
}

func (s *service) buildFilterForGetLanding(
	ctx			context.Context,
	name, tag 	string,
) (interface{}, error) {
	filter := bson.M{}

	if tag != "" {
		if err := s.buildFilterForTags(
			ctx,
			tag,
			&filter,
		); err != nil {
			return nil, err
		}
	}

	if name != "" {
		if err := s.buildNameFilterForGetLanding(
			ctx,
			name,
			&filter,
		); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (s *service) buildFilterForTags(
	ctx 		context.Context,
	tag			string,
	filter		*bson.M,	
) (error) {
	tags := strings.Split(tag, " ")

	ids, err := s.repository.GetIDsOfReposByLandingTags(
		ctx,
		tags,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	if len(ids) == 0 {
		return ErrTagNotFound
	}

	(map[string]interface{})(*filter)["repo_id"] = bson.M{"$in": ids}

	return nil
}

func (s *service) buildNameFilterForGetLanding(
	_		context.Context,
	name	string,
	filter	*bson.M,
) error {
	(map[string]interface{})(*filter)["title"] = bson.M{"$regex": name, "$options": "-i"}

	return nil
}

func (s *service) GetLanding(
	ctx context.Context, 
	ID uint64,
) (*landing.Landing, error) {
	logger := log.With(s.logger, "method", "GetLanding")
	l, err := s.repository.GetLandingByRepoID(
		ctx,
		ID,
	)
	if err == mongo.ErrNoDocuments {
		return nil, statuscode.WrapStatusError(
			ErrLandingNotFound,
			http.StatusNotFound,
		)
	} else if err != nil {
		level.Error(logger).Log("Err", err)
		return nil, statuscode.WrapStatusError(
			ErrFaieldToGetLanding,
			http.StatusInternalServerError,
		)
	}

	return l, nil
}