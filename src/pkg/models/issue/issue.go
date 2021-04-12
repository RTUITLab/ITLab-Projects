package issue

import (
	"github.com/ITLab-Projects/pkg/models/pullrequest"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/label"
	"github.com/ITLab-Projects/pkg/models/assignee"
	"github.com/ITLab-Projects/pkg/models/user"
)

type Issue struct {
	ID        			uint64     				`json:"id"`
	Number    			uint64     				`json:"number"`

	// GitLabNumber       	*uint64     			`json:"iid,omitempty"`

	Description 		string					`json:"body"`

	// GitlabDescription	string					`json:"description,omitempty"`

	Title     			string     				`json:"title"`
	User      			user.User       		`json:"user"`
	State     			string     				`json:"state"`
	Assignees 			[]assignee.Assignee 	`json:"assignees"`
	Milestone 			*milestone.Milestone	`json:"milestone,omitempty"`
	Labels				[]label.Label     		`json:"labels"`
	RepPath				string     				`json:"reppath"`
	ProjectPath			string     				`json:"project_path"`
	CreatedAt 			string     				`json:"created_at"`
	UpdatedAt 			string     				`json:"updated_at"`
	ClosedAt  			string     				`json:"closed_at"`
	HtmlUrl   			string     				`json:"html_url"`

	// GitLabHTMLUrl     	string 					`json:"web_url,omitempty"`
	
	PullRequest			pullrequest.PullRequest `json:"pull_request"`
}