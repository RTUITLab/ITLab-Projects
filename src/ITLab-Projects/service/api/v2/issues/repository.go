package issues

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
)

type Repository interface {
	IssuesRepository
	RepoRepository
	MilestoneRepository
}

type IssuesRepository interface {
	GetChunckedIssues(
		ctx 	context.Context,
		filter,
		sort	interface{},
		start,
		count 	int64,
	) ([]*milestone.IssuesWithMilestoneID, error)
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