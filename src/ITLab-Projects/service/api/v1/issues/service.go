package issues

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"golang.org/x/net/context"
)

type Service interface {
	GetIssues(
		ctx context.Context,
		start, 	count int64,
		name, 	tag string,
	) ([]*milestone.IssuesWithMilestoneID, error)

	GetLabels(
		ctx context.Context,
	) ([]interface{}, error)
}