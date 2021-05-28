package githubreq

import (
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/models/realese"
	"context"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/models/milestone"
)

type Requester interface {
	GetMilestonesForRepo(string) ([]milestone.Milestone, error)
	GetMilestonesForRepoWithID(repo.Repo) ([]milestone.MilestoneInRepo, error)
	GetRepositoriesWithoutURL() ([]repo.Repo, error)
	GetLastsRealeseWithRepoID(
		ctx context.Context,
		reps []repo.Repo,
		// error handling
		// if error is nil would'nt call
		f func(error),
	) ([]realese.RealeseInRepo, error)
	GetLastRealeseWithRepoID(rep repo.Repo) (realese.RealeseInRepo, error)
	GetRepositories() ([]repo.RepoWithURLS, error)
	GetAllMilestonesForRepoWithID(
		ctx context.Context,
		reps []repo.Repo,
		// error handling
		// if error is nil would'nt call
		f func(error),
	) ([]milestone.MilestoneInRepo, error)
	GetAllTagsForRepoWithID(
		ctx context.Context,
		reps []repo.Repo,
		// if f nill would'nt call
		f func(error),
	) ([]tag.Tag, error)
}