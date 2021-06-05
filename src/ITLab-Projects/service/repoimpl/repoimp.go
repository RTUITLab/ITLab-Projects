package repoimpl

import (
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/landing"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
)

type RepoImp struct {
	*estimate.EstimateRepositoryImp
	*issue.IssueRepositoryImp
	*functask.FuncTaskRepositoryImp
	*milestone.MilestoneRepositoryImp
	*reales.RealeseRepositoryImp
	*repo.RepoRepositoryImp
	*landing.LandingRepositoryImp
}

func New(
	Repo	*repositories.Repositories,
) *RepoImp {
	return &RepoImp{
		EstimateRepositoryImp: estimate.New(Repo.Estimate),
		IssueRepositoryImp: issue.New(Repo.Issue),
		FuncTaskRepositoryImp: functask.New(Repo.FuncTask),
		MilestoneRepositoryImp: milestone.New(Repo.Milestone),
		RealeseRepositoryImp: reales.New(Repo.Realese),
		RepoRepositoryImp: repo.New(Repo.Repo),
		LandingRepositoryImp: landing.New(Repo.Landing),
	}
}