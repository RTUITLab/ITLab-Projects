package updater

import (
	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"net/http"
	"context"
	"errors"
	"github.com/ITLab-Projects/pkg/statuscode"

	"github.com/go-kit/kit/log/level"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/updater"
	"github.com/go-kit/kit/log"
)

func init() {
	// to generate swagger
	_ = e.Message{}
}

var (
	ErrFailedToUpdateProjects	= errors.New("Failed to update projects")
)

type ServiceImp struct {
	repository		Repository
	logger 			log.Logger
	requester 		githubreq.Requester
	upd				*updater.Updater
}

func New(
	reposiotry 	Repository,
	logger		log.Logger,
	requester	githubreq.Requester,
	upd			*updater.Updater,
) *ServiceImp {
	return &ServiceImp{
		repository: reposiotry,
		logger: logger,
		requester: requester,
		upd: upd,
	}
}

// UpdateProjects
//
// @Tags projects
//
// @Summary Update all projects
//
// @Description make all request to github to update repositories, milestones
//
// @Security ApiKeyAuth
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
func (s *ServiceImp) UpdateProjects(
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

func (s *ServiceImp) resetUpdater() {
	if s.upd != nil {
		level.Debug(s.logger).Log("Reset update")
		s.upd.Reset()
	}
}

func (s *ServiceImp) getAllFromGithub(
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