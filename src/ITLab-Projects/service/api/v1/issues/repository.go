package issues

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
)

type Repository interface {
	IssueRepository
	RepoRepository
	MilestoneRepository
}

type IssueRepository interface {
	GetFiltrSortedFromToIssues(
		ctx context.Context,
		filter 	interface{},
		sort 	interface{},
		start 	int64,
		count 	int64,
	) ([]*milestone.IssuesWithMilestoneID, error)

	GetLabelsNameFromOpenIssues(
		ctx context.Context,
	) ([]interface{}, error)
}

type RepoRepository interface {
	GetReposAndScanTo(
		ctx 	context.Context,
		filter 	interface{},
		value 	interface{},
		options ...*options.FindOptions,
	) error
}

type MilestoneRepository interface {
	GetMilestonesAndScanTo(
		ctx context.Context,
		filter interface{},
		value interface{},
		opts ...*options.FindOptions,
	) error
}