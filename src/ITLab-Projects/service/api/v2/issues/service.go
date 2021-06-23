package issues

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
)

type Service interface {
	GetIssues(
		ctx 	context.Context,
		Query 	GetIssuesQuery,
	) ([]*milestone.IssuesWithMilestoneID, error)
}