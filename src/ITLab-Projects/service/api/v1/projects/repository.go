package projects

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/models/tag"
)

type Repository interface {
	RepoRepository
	MilestoneRepository
	IssueRepository
	RealeseRepository
	FuncTaskRepository
	EstimeateRepository
	LandingRepository
}

type RepoRepository interface{
	SaveReposAndSetDeletedUnfind(
		ctx context.Context,
		repos interface{},
	) error

	 GetFiltrSortFromToRepos(
		ctx 	context.Context,
		filter 	interface{},
		sort 	interface{},
		start 	int64,
		count 	int64,
	) ([]*repo.Repo, error)

	GetByID(
		ctx context.Context,
		ID uint64,
	) (*repo.Repo, error)

	DeleteByID(
		ctx context.Context,
		ID uint64,
	) error
}

type MilestoneRepository interface {
	SaveMilestonesAndSetDeletedUnfind(
		ctx context.Context,
		ms interface{},
	) error
	
	DeleteAllMilestonesByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) error


	GetAllMilestonesByRepoID(
		ctx 		context.Context,
		RepoID		uint64,
	) ([]*milestone.Milestone, error)
}

type IssueRepository interface {
	SaveIssuesAndSetDeletedUnfind(
		ctx context.Context,
		is 	interface{},
	) error

	GetAllIssuesByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) ([]*milestone.Issue, error)

	GetIssuesAndScanTo(
		ctx 		context.Context,
		filter 		interface{},
		value 		interface{},
		options 	...*options.FindOptions,
	) (error)

	DeleteAllIssuesByMilestonesID(
		ctx 		context.Context,
		MilestonesID []uint64,
	) error
}

type TagRepository interface {
	SaveAndDeleteUnfindTags(
		ctx context.Context,
		tgs interface{},
	) error

	DeleteTagsByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) (error)
	
	GetFilteredTags(
		ctx context.Context,
		filter interface{},
	) ([]*tag.Tag, error)

	GetFilteredTagsByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) ([]*tag.Tag, error)
	
}

type RealeseRepository interface {
	SaveRealeses(
		ctx context.Context,
		rs interface{},
	) error

	GetRealeseByRepoID(
		ctx 		context.Context,
		RepoID		uint64,
	) (*realese.RealeseInRepo, error)

	DeleteRealeseByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) error
}

type FuncTaskRepository interface {
	GetFuncTaskByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) (*functask.FuncTaskFile, error)

	GetFuncTasksByMilestonesID(
		ctx 			context.Context,
		MilestonesID	[]uint64,
	) ([]*functask.FuncTaskFile, error)

	DeleteManyFuncTasksByMilestonesID(
		ctx 			context.Context,
		MilestonesID	[]uint64,
	) error
}

type EstimeateRepository interface {
	GetEstimateByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) (*estimate.EstimateFile, error)

	GetEstimatesByMilestonesID(
		ctx 			context.Context,
		MilestonesID	[]uint64,
	) ([]*estimate.EstimateFile, error)

	DeleteManyEstimatesByMilestonesID(
		ctx 			context.Context,
		MilestonesID	[]uint64,
	) error
}

type LandingRepository interface {
	SaveAndDeleteUnfindLanding(
		ctx context.Context,
		ls interface{},
	) error

	DeleteLandingsByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) (error)
	
	GetFilteredLandings(
		ctx context.Context,
		filter interface{},
	) ([]*landing.Landing, error)

	GetFilteredLandingsByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) ([]*landing.Landing, error)

	GetIDsOfReposByLandingTags(
		ctx		context.Context,
		Tags	[]string,
	) ([]uint64, error)

	GetLandingTagsByRepoID(
		ctx		context.Context,
		RepoID	uint64,
	) ([]*tag.Tag, error)

}