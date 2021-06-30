package issues

import (
	"context"

	"github.com/ITLab-Projects/pkg/models/milestone"
)

type Service interface {
	GetIssues(
		ctx 	context.Context,
		Query	GetIssuesQuery,
	) ([]*milestone.IssuesWithMilestoneID, error)

	GetLabels(
		ctx context.Context,
	) ([]interface{}, error)
}