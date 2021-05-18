package projects

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

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
) (*repoasproj.RepoAsProj, error) {
	rep, err := s.repository.GetByID(
		ctx,
		ID,
	)
	if err != nil {
		return nil, err
	}

	
}

func (s *service) GetProjects(
	ctx 			context.Context, 
	start, 	count 	int64,
	name, 	tag		string,
) ([]*repoasproj.RepoAsProjCompactPointers, error) {
	filter, err := s.buildFilterForGetProject(
		ctx,
		name,
		tag,
	)
	if err != nil {
		return nil, err
	}

	repos, err := s.repository.GetFiltrSortFromToRepos(
		ctx,
		filter,
		bson.D{ {"createdat", -1}, {"deleted", 1}},
		start,
		count,
	)
	if err != nil {
		return nil, err
	}

	projs, err := s.GetCompatcProj(ctx, repos)
	if err != nil {
		return nil, err
	}
	sort.Sort(repoasproj.ByCreateDatePointers(projs))
	return projs, nil
}

func (s *service) UpdateProjects(
	ctx context.Context,
) error {
	s.resetUpdater()
	defer s.resetUpdater()

	repos, ms, rs, tgs, err := s.getAllFromGithub(ctx)
	if err != nil {
		return err
	}

	if err := s.repository.SaveReposAndSetDeletedUnfind(
		ctx,
		repos,
	); err != nil {
		return err
	}

	if err := s.repository.SaveMilestonesAndSetDeletedUnfind(
		ctx,
		ms,
	); err != nil {
		return err
	}

	is := getIssuesFromMilestone(ms)

	if err := s.repository.SaveIssuesAndSetDeletedUnfind(
		ctx,
		is,
	); err != nil {
		return err
	}

	if err := s.repository.SaveRealeses(
		ctx,
		rs,
	); err != nil {
		return err
	}

	if err := s.repository.SaveAndDeleteUnfindTags(
		ctx,
		tgs,
	); err != nil {
		return nil
	}

	return nil
}

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
	if err != nil {
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
	if err!= nil {
		return nil, err
	}

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

		is, err := s.repository.Get
	}
}