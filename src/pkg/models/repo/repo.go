package repo

import (
	. "github.com/ITLab-Projects/pkg/models/project"
	. "github.com/ITLab-Projects/pkg/models/user"
)

type Repo struct {
	ID          		uint64 		`json:"id"`
	Platform			string		`json:"platform,omitempty"`
	Name        		string 		`json:"name"`
	Contributors		[]User		`json:"contributors"`
	Path				string		`json:"path_with_namespace,omitempty"`
	HTMLUrl     		string 		`json:"html_url"`

	// GitLabHTMLUrl     	string 		`json:"web_url,omitempty"`

	Description 		string 		`json:"description"`
	CreatedAt   		string 		`json:"created_at"`
	UpdatedAt   		string 		`json:"updated_at"`
	PushedAt			string		`json:"pushed_at"`

	// GitLabUpdatedAt   	string 		`json:"last_activity_at,omitempty"`

	Language    		string 		`json:"language"`
	Languages			map[string]int	`json:"languages"`
	Archived    		bool   		`json:"archived"`
	OpenIssues  		int			`json:"open_issues_count"`
	Meta				Meta		`json:"meta"`
}