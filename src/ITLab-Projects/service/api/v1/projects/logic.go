package projects

import (
	"fmt"
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/ITLab-Projects/service/api/v1/beforedelete"

	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
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
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/updater"
	"github.com/go-kit/kit/log"
)

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

func (s *service) GetProjects(
	ctx 			context.Context, 
	start, 	count 	int64,
	name, 	tag		string,
) ([]*repoasproj.RepoAsProjCompactPointers, error) {
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
		return nil, err
	}
	sort.Sort(repoasproj.ByCreateDatePointers(projs))
	return projs, nil
}

func (s *service) UpdateProjects(
	ctx context.Context,
) error {
	logger := log.With(s.logger, "method", "UpdateProjects")
	s.resetUpdater()
	defer s.resetUpdater()

	repos, ms, rs, tgs, err := s.getAllFromGithub(ctx)
	switch {
	case err == githubreq.ErrForbiden, err == githubreq.ErrUnatorizared:
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects, //TODO think to make it real error to catch in future
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

	if err := s.repository.SaveAndDeleteUnfindTags(
		ctx,
		tgs,
	); err != nil {
		level.Error(logger).Log("Failed to update projects: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToUpdateProjects,
			http.StatusInternalServerError,
		)
	}

	return nil
}

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

	if err := s.repository.DeleteTagsByRepoID(
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

	if err := s.repository.DeleteTagsByRepoID(
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
) ([]repo.Repo, []milestone.MilestoneInRepo, []realese.RealeseInRepo, []tag.Tag, error) {
	repos, err := s.requester.GetRepositories()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	msChan := make(chan []milestone.MilestoneInRepo, 1)
	rsChan := make(chan []realese.RealeseInRepo, 1)
	tgsChan := make(chan []tag.Tag, 1)
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
		tgs, err := s.requester.GetAllTagsForRepoWithID(
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

		tgsChan <- tgs
	}()

	var (
		ms 	[]milestone.MilestoneInRepo = nil
		rs 	[]realese.RealeseInRepo		= nil
		tgs	[]tag.Tag					= nil
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
		case _tgs := <- tgsChan:
			tgs = _tgs
		}
	}

	close(msChan)
	close(rsChan)
	close(tgsChan)

	return repo.ToRepo(repos) ,ms, rs, tgs, nil
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

	tags, err := s.repository.GetFilteredTags(
		ctx,
		bson.M{"tag": bson.M{"$in": massOfTags}},
	)
	if err == mongo.ErrNoDocuments {
		return nil
	}else if err != nil {
		return err
	}

	var ids []uint64
	for _, t := range tags {
		ids = append(ids, t.RepoID)
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

			tgs, err := s.repository.GetFilteredTagsByRepoID(
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

	tgs, err := s.repository.GetFilteredTagsByRepoID(
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