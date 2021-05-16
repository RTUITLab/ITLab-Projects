package projects

import (
	"errors"
	"github.com/ITLab-Projects/pkg/models/repo"
	"context"

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

type service struct {
	repository 		Repository
	logger 			log.Logger
	requester		githubreq.Requester
	mfsRequester	mfsreq.Requester
	upd				*updater.Updater
}


func (s *service) GetProject(
	ctx context.Context,
	ID uint64,
) (*repoasproj.RepoAsProj, error) {
	panic("not implemented") // TODO: Implement
}

func (s *service) GetProjects(
	ctx context.Context, 
	start int64, 
	count int64,
) ([]*repoasproj.RepoAsProjCompact, error) {
	panic("not implemented") // TODO: Implement
}

// func (s *service) UpdateProjects(
// 	ctx context.Context,
// ) error {
// 	s.resetUpdater()
// 	defer s.resetUpdater()

// 	repos, ms, rs, tgs, err := s.getAllFromGithub(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if err := s.repository.SaveReposAndSetDeletedUnfind(
// 		ctx,
// 		repos,
// 	); err != nil {
// 		return err
// 	}

// 	if err := s.repository.SaveMilestonesAndSetDeletedUnfind(
// 		ctx,
// 		ms,
// 	); err != nil {
// 		return err
// 	}


// }

func (s *service) DeleteProject(
	ctx context.Context, 
	ID uint64,
) error {
	panic("not implemented") // TODO: Implement
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
