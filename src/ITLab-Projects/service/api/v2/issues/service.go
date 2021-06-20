package issues

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
)

type Service interface {
	GetIssues(
		ctx context.Context,
		start, 	count int64,
		name, 	tag string,
	) ([]*milestone.IssuesWithMilestoneID, error)
}