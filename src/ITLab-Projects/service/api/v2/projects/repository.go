package projects

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/models/repo"
	"context"
)

type Repository interface {
	RepoRepository
	LandingRepository
	MilestoneRepository
}

type RepoRepository interface {
	GetChunckedRepos(
		ctx 	context.Context,
		filter,
		sort	interface{},
		start,
		count 	int64,
	) ([]*repo.Repo, error)
}

type LandingRepository interface {
	GetIDsOfReposByLandingTags(
		ctx		context.Context,
		Tags	[]string,
	) ([]uint64, error)

	GetLandingTagsByRepoID(
		ctx		context.Context,
		RepoID	uint64,
	) ([]*tag.Tag, error)
}

type MilestoneRepository interface {
	GetAllMilestonesByRepoID(
		ctx 		context.Context,
		RepoID		uint64,
	) ([]*milestone.Milestone, error)
}
