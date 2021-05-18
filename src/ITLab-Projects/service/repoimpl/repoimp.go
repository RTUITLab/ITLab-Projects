package repoimpl

import (
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/ITLab-Projects/service/repoimpl/tag"
)

type RepoImp struct {
	*estimate.EstimateRepositoryImp
	*issue.IssueRepositoryImp
	*functask.FuncTaskRepositoryImp
	*milestone.MilestoneRepositoryImp
	*reales.RealeseRepositoryImp
	*repo.RepoRepositoryImp
	*tag.TagRepositoryImp
}