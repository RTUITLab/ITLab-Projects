package milestone

import (
	"github.com/ITLab-Projects/pkg/models/assignee"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/label"
	"github.com/ITLab-Projects/pkg/models/pullrequest"
	"github.com/ITLab-Projects/pkg/models/user"
	"github.com/Kamva/mgm"
)

type MilestoneFromGH struct {
	ID        			uint64     	`json:"id"`
	Number    			uint64     	`json:"number"`
	State     			string     	`json:"state"`
	Title     			string     	`json:"title"`
	Description 		string 		`json:"description"`
	Creator      		user.User   `json:"creator"`
	OpenIssues			uint64		`json:"open_issues"`
	ClosedIssues		uint64		`json:"closed_issues"`
}

type Milestone struct {
	MilestoneFromGH						`bson:",inline"`
	Issues 			[]Issue				`json:"issues" bson:"-"`
	Estimate 		*estimate.Estimate 	`json:"estimate" bson:"-"`
	FuncTask		*functask.FuncTask	`json:"func_task" bson:"-"`
	Deleted			bool				`json:"deleted" bson:"deleted"`
}

type MilestoneInRepo struct {
	mgm.DefaultModel				`json:"-" bson:",inline"`
	RepoID				uint64		`json:"repo_id"`
	Milestone						`bson:",inline"`
}

func (m *MilestoneInRepo) CollectionName() string {
	return "milestones"
}

func GetIDS(ms []MilestoneInRepo) []uint64 {
	var ids []uint64
	for _, m := range ms {
		ids = append(ids, m.Milestone.ID)
	}

	return ids
}

type Issue struct {
	ID        			uint64     				`json:"id"`
	Number    			uint64     				`json:"number"`
	Description 		string					`json:"body"`
	Title     			string     				`json:"title"`
	User      			user.User       		`json:"user"`
	State     			string     				`json:"state"`
	Assignees 			[]assignee.Assignee 	`json:"assignees"`
	Labels				[]label.Label     		`json:"labels"`
	RepPath				string     				`json:"reppath"`
	ProjectPath			string     				`json:"project_path"`
	CreatedAt 			string     				`json:"created_at"`
	UpdatedAt 			string     				`json:"updated_at"`
	ClosedAt  			string     				`json:"closed_at"`
	HtmlUrl   			string     				`json:"html_url"`
	PullRequest			pullrequest.PullRequest `json:"pull_request"`
}

type IssueFromGH struct {
	Issue
	Milestone *MilestoneFromGH 				`json:"milestone,omitempty"`
}

type IssuesWithMilestoneID struct {
	mgm.DefaultModel				`bson:",inline"`
	MilestoneID			uint64		`json:"milestone_id" bson:"milestone_id"`
	RepoID				uint64		`json:"repo_id" bson:"repo_id"`
	Issue							`bson:",inline"`
	Deleted				bool		`json:"deleted" bson:"deleted"`
}

func (i *IssuesWithMilestoneID) CollectionName() string {
	return "issues"
}