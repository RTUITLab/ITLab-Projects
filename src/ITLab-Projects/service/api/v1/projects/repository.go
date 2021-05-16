package projects

import (
	"github.com/ITLab-Projects/pkg/models/repo"
	"context"
)

type Repository interface {
	RepoRepository
	MilestoneRepository
	IssueRepository
	TagRepository
	RealeseRepository
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
	
	DeleteAllByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) error

	DeleteAllByMilestoneID(
		ctx 		context.Context,
		MilestoneID uint64,
	) error
}

type IssueRepository interface {
	SaveIssuesAndSetDeletedUnfind(
		ctx context.Context,
		is 	interface{},
	) error

	DeleteAllTagsByRepoID(
		ctx 		context.Context,
		RepoID uint64,
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

	
}

type RealeseRepository interface {
	SaveRealeses(
		ctx context.Context,
		rs interface{},
	) error
}