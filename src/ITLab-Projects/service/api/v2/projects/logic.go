package projects

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repo"
	"strings"
	"sort"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"github.com/ITLab-Projects/pkg/statuscode"
	"errors"
	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"context"
	"github.com/go-kit/kit/log"
)

func init() {
	// to generate swagger
 	_ = e.Message{}
}

var (
	ErrTagNotFound				= errors.New("Tag not found")
	ErrFailedToGetProjects		= errors.New("Failed to get projects")
)

type ServiceImp struct {
	logger log.Logger
	repository Repository
}

func New(
	repository Repository,
	logger log.Logger,
) *ServiceImp {
	return &ServiceImp{
		logger: logger,
		repository: repository,
	}
}

// GetProjects
// 
// @Tags projects
// 
// @Summary return projects according to query value
// 
// @Description return a projects you can filter count of them
// 
// @Description tags, name
// 
// @Produce json
// 
// @Security ApiKeyAuth
// 
// @Router /v2/projects [get]
// 
// @Param start query integer false "represents the number of skiped projects"
// 
// @Param count query integer false "represent a limit of projects, standart and max count equal 50"
// 
// @Param tag query string false "use to filter projects by tag"
// 
// @Param name query string false "use to filter by name"
// 
// @Success 200 {object} projects.GetProjectsResp
// 
// @Failure 500 {object} e.Message "Failed to get projects"
// 
// @Failure 500 {object} e.Message "Failed to get repositories"
// 
// @Failure 401 {object} e.Message 
func (s *ServiceImp) GetProjects(
	ctx 			context.Context,
	start, 	count 	int64,
	name, 	tag		string,
) ([]*repoasproj.RepoAsProjCompactPointers, error) {
	if count == 0 || count > 50 {
		count = 50
	}
	logger := log.With(s.logger, "method", "GetProjects")
	filter, err := s.buildFilterForGetProject(
		ctx,
		name,
		tag,
	)
	if err == ErrTagNotFound {
		return []*repoasproj.RepoAsProjCompactPointers{}, nil
	} else if err != nil {
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetProjects,
			http.StatusInternalServerError,
		)
	}

	repos, err := s.repository.GetChunckedRepos(
		ctx,
		filter,
		bson.D{ {"createdat", -1}, {"deleted", 1}},
		start,
		count,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		level.Info(logger).Log("Not found projects for this filters: filter", filter)
		return []*repoasproj.RepoAsProjCompactPointers{}, nil
	case err != nil:
		level.Error(logger).Log("Failed to get projects: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetProjects,
			http.StatusInternalServerError,
		)
	}

	projs, err := s.GetCompatcProj(ctx, repos)
	if err != nil {
		level.Error(logger).Log("Failed to get projects: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetProjects,
			http.StatusInternalServerError,
		)
	}
	sort.Sort(repoasproj.ByCreateDatePointers(projs))
	return projs, nil
}

func (s *ServiceImp) buildFilterForGetProject(
	ctx 		context.Context,
	name, tag 	string,
) (interface{}, error) {
	filter := bson.M{}

	if name != "" {
		if err := s.buildNameFilterForGetProjects(
			ctx,
			name,
			&filter,
		); err != nil {
			return nil, err
		}
	}

	if tag != "" {
		if err := s.buildTagFilterForGetProjects(
			ctx,
			tag,
			&filter,
		); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (s *ServiceImp) buildTagFilterForGetProjects(
	ctx 	context.Context,
	t 		string,
	filter	*bson.M,
) (error) {
	massOfTags := strings.Split(t, " ")

	ids, err := s.repository.GetIDsOfReposByLandingTags(
		ctx,
		massOfTags,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}
	if len(ids) > 0 {
		(map[string]interface{})(*filter)["id"] = bson.M{"$in": ids}
	} else {
		return ErrTagNotFound
	}

	return nil
}

func (s *ServiceImp) buildNameFilterForGetProjects(
	_ context.Context,
	name string,
	filter *bson.M,
) error {
	(map[string]interface{})(*filter)["name"] = bson.M{"$regex": name, "$options": "-i"}

	return nil
}

func (s *ServiceImp) GetCompatcProj(
	ctx context.Context,
	repos []*repo.Repo,
) ([]*repoasproj.RepoAsProjCompactPointers, error) {
	ctx, cancel := context.WithCancel(
		ctx,
	)
	defer cancel()

	errChan := make(chan error, 1)
	projChan := make(chan *repoasproj.RepoAsProjCompactPointers, 1)

	var count uint
	for i := range repos {
		count++
		go func(r *repo.Repo) {
			defer s.catchPanic()
			proj := &repoasproj.RepoAsProjCompactPointers{
				Repo: r,
			}

			if ms, err := s.repository.GetAllMilestonesByRepoID(
				ctx,
				r.ID,
			); err == mongo.ErrNoDocuments {
				// Pass
			} else if err != nil {
				cancel()
				errChan <- err
				return
			} else {
				proj.Completed = countCompleted(ms)
			}

			tgs, err := s.repository.GetLandingTagsByRepoID(
				ctx,
				r.ID,
			)
			if err == mongo.ErrNoDocuments {
				// Pass
			} else if err != nil {
				cancel()
				errChan <- err
				return
			} else {
				proj.Tags = tgs
			}

			projChan <- proj
		}(repos[i])
	}

	var projs []*repoasproj.RepoAsProjCompactPointers

	for i := uint(0); i < count; i++ {
		select{
		case p := <- projChan:
			projs = append(projs, p)
		case <-ctx.Done():
			err := <- errChan
			defer close(errChan)
			defer close(projChan)
			return nil, err
		}
	}

	return projs, nil
}

func (s *ServiceImp) catchPanic() {
	if r := recover(); r != nil {
		level.Debug(s.logger).Log("CatchPanic", r)
	}
}

func countCompleted(ms []*milestone.Milestone) float64 {
	var open float64
	var closed float64
	for _, m := range ms {
		if m.OpenIssues != 0 {
			open += float64(m.OpenIssues)
			closed += float64(m.ClosedIssues)
		}
	}

	if open + closed == 0 {
		return 1
	} else {
		return (closed)/(open+closed)
	}
}