package milestone

import (
	"github.com/ITLab-Projects/pkg/models/pullrequest"
	"github.com/ITLab-Projects/pkg/models/label"
	"github.com/ITLab-Projects/pkg/models/assignee"
	"github.com/ITLab-Projects/pkg/models/user"
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
	MilestoneFromGH			`bson:",inline"`
	Issues 			[]Issue	`json:"issues"`
}

type MilestoneInRepo struct {
	RepoID				uint64		`json:"repo_id"`
	Milestone						`bson:",inline"`
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