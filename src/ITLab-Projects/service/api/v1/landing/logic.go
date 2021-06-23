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
// GetAllLanding
// 
// @Summary return all landings according to path params
// 
// @Tags projects
// 
// @Produce json
// 
// @Router /v1/projects/landing [get]
// 
// @Param start query integer false "represent how much landins need to skip"
// 
// @Param count query integer false "represent a max count of returing landing"
// 
// @Param tag query string false "return a landings with this tags"
// 
// @Param name query string false "return landing with this names"
// 
// @Success 200 {array} landing.LandingCompact
// 
// @Failure 500 {object} e.Message
func (s *service) GetAllLandings(
	ctx 	context.Context, 
	Query	GetAllLandingsQuery,
) ([]*landing.LandingCompact, error) {
	logger := log.With(s.logger, "method", "GetAllLandings")
	filter, err := s.buildFilterForGetLanding(
		ctx,
		Query.Name,
		Query.Tag,
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
		Query.Start,
		Query.Count,
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

// GetLanding
// 
// @Tags projects
// 
// @Summary return a current landing
// 
// @Description return a landing according to id
// 
// @Produce json
// 
// @Router /v1/projects/landing/{id} [get]
// 
// @Param id path integer true "id of landing"
// 
// @Success 200 {object} landing.Landing
// 
// @Failure 404 {object} e.Message
// 
// @Failure 500 {object} e.Message
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