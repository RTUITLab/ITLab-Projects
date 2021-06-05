package projects

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	e "github.com/ITLab-Projects/pkg/err"

	"github.com/ITLab-Projects/service/api/v1/beforedelete"

	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/statuscode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-kit/kit/log/level"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/ITLab-Projects/pkg/updater"
	"github.com/go-kit/kit/log"
)

func init() {
	// to generate swagger
	_ = e.Message{}
}

var (
	ErrProjectNotFound 			= errors.New("Poject not found")
	ErrFaieldToGetProject		= errors.New("Failed to get project")
	ErrFailedToGetProjects		= errors.New("Failed to get projects")
	ErrFailedToUpdateProjects	= errors.New("Failed to update projects")
	ErrFailedToDeleteProject	= errors.New("Failed to delete project")
)

type service struct {
	repository 		Repository
	logger 			log.Logger
	requester		githubreq.Requester
	mfsRequester	mfsreq.Requester
	upd				*updater.Updater
}

func New(
	repository 	Repository,
	logger 		log.Logger,
	requester	githubreq.Requester,
	mfsreq		mfsreq.Requester,
	upd			*updater.Updater,
) *service {
	return &service{
		repository: repository,
		logger: logger,
		requester: requester,
		mfsRequester: mfsreq,
		upd: upd,
	}
}

// GetProject
// 
// @Summary return project according to id
// 
// @Tags projects
// 
// @Description return a project according to id value in path
// 
// @Produce json
// 
// @Router /v1/projects/{id} [get]
// 
// @Param id path integer true "a uint value of repository id"
// 
// @Success 200 {object} repoasproj.RepoAsProj
// 
// @Failure 404 {object} e.Message
// 
// @Failure 500 {object} e.Message
// 
// @Failure 401 {object} e.Message 
func (s *service) GetProject(
	ctx context.Context,
	ID uint64,
) (*repoasproj.RepoAsProjPointer, error) {
	logger := log.With(s.logger, "method", "GetProject")
	rep, err := s.repository.GetByID(
		ctx,
		ID,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, statuscode.WrapStatusError(
			ErrProjectNotFound,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to get project: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFaieldToGetProject,
			http.StatusInternalServerError,
		)
	}

	proj, err := s.getProjects(
		ctx,
		rep,
	)
	if err != nil {
		level.Error(logger).Log("Failed to get project: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFaieldToGetProject,
			http.StatusInternalServerError,
		)
	}

	return proj, nil
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
// @Router /v1/projects [get]
// 
// @Param start query integer false "represents the number of skiped projects"
// 
// @Param count query integer false "represent a limit of projects, standart and max count equal 50"
// 
// @Param tag query string false "use to filter projects by tag"
// 
// @Param name query string false "use to filter by name"
// 
// @Success 200 {array} repoasproj.RepoAsProjCompact
// 
// @Failure 500 {object} e.Message "Failed to get projects"
// 
// @Failure 500 {object} e.Message "Failed to get repositories"
// 
// @Failure 401 {object} e.Message 
func (s *service) GetProjects(
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
	if err != nil {
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetProjects,
			http.StatusInternalServerError,
		)
	}

	repos, err := s.repository.GetFiltrSortFromToRepos(
		ctx,
		filter,
		bson.D{ {"createdat", -1}, {"deleted", 1}},
		start,
		count,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		level.Info(logger).Log("Not found projects for this filters: filter", filter)
		return nil, nil
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


// UpdateProjects
//
// @Tags projects
//
// @Summary Update all projects
//
// @Description make all request to github to update repositories, milestones
//
// @Router /v1/projects [post]
//
// @Success 200
//
// @Failure 409 {object} e.Err
//
// @Failure 500 {object} e.Message
//
// @Failure 401 {object} e.Message
//
// @Failure 403 {object} e.Message "if you are nor admin"
func (s *service) UpdateProjects(
	ctx context.Context,
) error {
	logger := log.With(s.logger, "method", "UpdateProjects")
	if !updater.IsUpdateContext(ctx) {
		level.Debug(logger).Log("From user")
		s.resetUpdater()
		defer s.resetUpdater()
	}

	repos, ms, rs, ls, err := s.getAllFromGithub(ctx)
	switch {
	case err == githubreq.ErrForbiden, err == githubreq.ErrUnatorizared:
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusConflict,
		)
	case err != nil:
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveReposAndSetDeletedUnfind(
		ctx,
		repos,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveMilestonesAndSetDeletedUnfind(
		ctx,
		ms,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	is := getIssuesFromMilestone(ms)

	if err := s.repository.SaveIssuesAndSetDeletedUnfind(
		ctx,
		is,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveRealeses(
		ctx,
		rs,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveAndDeleteUnfindLanding(
		ctx,
		ls,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	return nil
}

// DeleteProject
// 
// @Summary delete project by id
// 
// @Description delete project by id and all milestones in it
// 
// @Tags projects
// 
// @Router /v1/projects/{id} [delete]
// 
// @Param id path integer true "id of project"
// 
// @Success 200
// 
// @Failure 404 {object} e.Message
// 
// @Failure 500 {object} e.Message
// 
// @Failure 409 {object} e.Message "some problems with microfileservice"
// 
// @Failure 401 {object} e.Message 
// 
// @Failure 403 {object} e.Message "if you are not admin"
func (s *service) DeleteProject(
	ctx context.Context, 
	ID 	uint64,
	// For mfs requester
	r 	*http.Request,
) error {
	logger := log.With(s.logger,"method", "DeleteProjects")
	rep, err := s.repository.GetByID(
		ctx,
		ID,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return statuscode.WrapStatusError(
			ErrProjectNotFound,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteProject,
			http.StatusInternalServerError,
		)
	}

	err = s.deleteMilestone(
		ctx,
		rep.ID,
		beforedelete.BeforeDeleteWithReq(
			s.mfsRequester,
			r,
		),
	)
	switch {
	case errors.Is(err, mfsreq.NetError):
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			mfsreq.NetError,
			http.StatusConflict,
		)
	case mfsreq.IfUnexcpectedCode(err):
		uce := err.(*mfsreq.UnexpectedCodeErr)
		causedErr := fmt.Errorf("Unecxpected code from microfileserver: %v", uce.Code)
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			causedErr,
			http.StatusConflict,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteProject,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.DeleteLandingsByRepoID(
		ctx,
		rep.ID,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteProject,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.DeleteRealeseByRepoID(
		ctx,
		rep.ID,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteProject,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.DeleteByID(
		ctx,
		rep.ID,
	); err != nil {
		level.Error(logger).Log("Failed to delete project: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteProject,
			http.StatusInternalServerError,
		)
	}

	return nil
}

func (s *service) resetUpdater() {
	if s.upd != nil {
		level.Debug(s.logger).Log("Reset update")
		s.upd.Reset()
	}
}

func (s *service) getAllFromGithub(
	ctx context.Context,
) ([]repo.Repo, []milestone.MilestoneInRepo, []realese.RealeseInRepo, []*landing.Landing, error) {
	repos, err := s.requester.GetRepositories()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	msChan := make(chan []milestone.MilestoneInRepo, 1)
	rsChan := make(chan []realese.RealeseInRepo, 1)
	lsChan := make(chan []*landing.Landing, 1)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(
		ctx,
	)
	defer cancel()
	logger := log.With(s.logger, "method", "getAllFromGithub")

	go func() {
		ms, err := s.requester.GetAllMilestonesForRepoWithID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					cancel()
					errChan <- err
				} else if errors.Unwrap(err) == githubreq.UnexpectedCode {
					// Pass
				} else {
					level.Warn(logger).Log(err)
				}
			},
		)
		if err != nil {
			return
		}

		msChan <- ms
	}()

	go func() {
		rs, err := s.requester.GetLastsRealeseWithRepoID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					cancel()
					errChan <- err
				} else if errors.Unwrap(err) == githubreq.UnexpectedCode {
					// Pass
				} else {
					level.Warn(logger).Log(err)
				}
			},
		)
		if err != nil {
			return
		}
		rsChan <- rs
	}()

	go func() {
		ls, err := s.requester.GetAllLandingsForRepoWithID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					cancel()
					errChan <- err
				} else if errors.Unwrap(err) == githubreq.UnexpectedCode {
					// Pass
				} else { 
					level.Warn(logger).Log(err)
				}
			},
		)
		if err != nil {
			return
		}

		lsChan <- ls
	}()

	var (
		ms 	[]milestone.MilestoneInRepo = nil
		rs 	[]realese.RealeseInRepo		= nil
		ls	[]*landing.Landing			= nil
	)


	for i := 0; i < 3; i++ {
		select {
		case <- ctx.Done():
			err := <- errChan
			return nil, nil, nil, nil, err
		case _ms := <- msChan:
			ms = _ms
		case _rs := <- rsChan:
			rs = _rs
		case _ls := <- lsChan:
			ls = _ls
		}
	}

	close(msChan)
	close(rsChan)
	close(lsChan)

	return repo.ToRepo(repos) ,ms, rs, ls, nil
}

func getIssuesFromMilestone(
	ms []milestone.MilestoneInRepo,
) []milestone.IssuesWithMilestoneID {
	var is []milestone.IssuesWithMilestoneID

	isChan := make(chan milestone.IssuesWithMilestoneID)

	count := 0
	for j, _ := range ms {
		for i, _ := range ms[j].Issues {
			count++
			go func(i milestone.Issue, MID , RepoID uint64) {
				isChan <- milestone.IssuesWithMilestoneID{
					MilestoneID: MID,
					RepoID: RepoID,
					Issue: i,
				}
			}(ms[j].Issues[i], ms[j].Milestone.ID, ms[j].RepoID)
		}
	}

	for i := 0; i < count; i++ {
		select{
		case issue := <- isChan:
			is = append(is, issue)
		}
	}

	return is
}

func (s *service) buildFilterForGetProject(
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

func (s *service) buildTagFilterForGetProjects(
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

	(map[string]interface{})(*filter)["id"] = bson.M{"$in": ids}

	return nil
}

func (s *service) buildNameFilterForGetProjects(
	_ context.Context,
	name string,
	filter *bson.M,
) error {
	(map[string]interface{})(*filter)["name"] = bson.M{"$regex": name, "$options": "-i"}

	return nil
}

func (s *service) GetCompatcProj(
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

func (s *service) catchPanic() {
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

func (s *service) getProjects(
	ctx context.Context,
	rep *repo.Repo,
) (*repoasproj.RepoAsProjPointer, error) {
	proj := &repoasproj.RepoAsProjPointer{
		Repo: rep,
	}

	ms, err := s.repository.GetAllMilestonesByRepoID(
		ctx,
		rep.ID,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}
	proj.Milestones = ms

	if err := s.getAssetsForMilestones(
		ctx,
		ms,
	); err != nil {
		return nil, err
	}

	proj.Completed = countCompleted(ms)

	if rls, err := s.repository.GetRealeseByRepoID(
		ctx,
		rep.ID,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	} else {
		proj.LastRealese = &rls.Realese
	}

	tgs, err := s.repository.GetLandingTagsByRepoID(
		ctx,
		rep.ID,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	} else {
		proj.Tags = tgs
	}

	return proj, nil
}

func (s *service) getAssetsForMilestones(
	ctx context.Context,
	ms []*milestone.Milestone,
) error {
	for i, _ := range ms {
		est, err := s.repository.GetEstimateByMilestoneID(
			ctx,
			ms[i].ID,
		)
		if err == mongo.ErrNoDocuments {
			// Pass
		} else if err != nil {
			return err
		} else {
			ms[i].Estimate = &estimate.Estimate{
				MilestoneID: est.MilestoneID,
				EstimateURL: s.mfsRequester.GenerateDownloadLink(est.FileID),
			}
		}

		f, err := s.repository.GetFuncTaskByMilestoneID(
			ctx,
			ms[i].ID,
		)
		if err == mongo.ErrNoDocuments {
			// Pass
		} else if err != nil {
			return err
		} else {
			ms[i].FuncTask = &functask.FuncTask{
				MilestoneID: f.MilestoneID,
				FuncTaskURL: s.mfsRequester.GenerateDownloadLink(f.FileID),
			}
		}

		if err := s.repository.GetIssuesAndScanTo(
			ctx,
			bson.M{"milestone_id": ms[i].ID},
			&ms[i].Issues,
			options.Find(),
		); err == mongo.ErrNoDocuments {
			// Pass
			level.Debug(s.logger).Log("Don't finc issues")
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) deleteMilestone(
	ctx 		context.Context,
	RepoID		uint64,
	f			beforedelete.BeforeDeleteFunc,
) error {
	ms, err := s.repository.GetAllMilestonesByRepoID(
		ctx,
		RepoID,
	)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}

	MilestonesID := func(ms []*milestone.Milestone) []uint64 {
		var ids []uint64
		for _, m := range ms {
			ids = append(ids, m.ID)
		}
		return ids
	}(ms)

	if err := s.deleteEstimates(
		ctx,
		MilestonesID,
		f,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	if err := s.deleteFuncTask(
		ctx,
		MilestonesID,
		f,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	if err := s.repository.DeleteAllIssuesByMilestonesID(
		ctx,
		MilestonesID,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	if err := s.repository.DeleteAllMilestonesByRepoID(
		ctx,
		RepoID,
	); err != nil {
		return err
	}

	return nil
}

func (s *service) deleteEstimates(
	ctx 			context.Context,
	MilestonesID	[]uint64,
	f				beforedelete.BeforeDeleteFunc,
) error {
	est, err := s.repository.GetEstimatesByMilestonesID(
		ctx,
		MilestonesID,
	)
	if err != nil {
		return err
	}
	
	if err := f(est); err != nil {
		return err
	}

	if err := s.repository.DeleteManyEstimatesByMilestonesID(
		ctx,
		MilestonesID,
	); err != nil {
		return err
	}

	return nil
}

func (s *service) deleteFuncTask(
	ctx 			context.Context,
	MilestonesID 	[]uint64,
	f				beforedelete.BeforeDeleteFunc,
) error {
	tasks, err := s.repository.GetFuncTasksByMilestonesID(
		ctx,
		MilestonesID,
	)
	if err != nil {
		return err
	}

	if err := f(tasks); err != nil {
		return err
	}

	if err := s.repository.DeleteManyFuncTasksByMilestonesID(
		ctx,
		MilestonesID,
	); err != nil {
		return err
	}

	return nil
}