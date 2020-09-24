package models

type Repos struct {
	ID          		uint64 		`json:"id"`
	Platform			string		`json:"platform,omitempty"`
	Name        		string 		`json:"name"`
	Contributors		[]User		`json:"contributors"`
	Path				string		`json:"path_with_namespace,omitempty"`
	HTMLUrl     		string 		`json:"html_url"`
	GitLabHTMLUrl     	string 		`json:"web_url,omitempty"`
	Description 		string 		`json:"description"`
	CreatedAt   		string 		`json:"created_at"`
	UpdatedAt   		string 		`json:"updated_at"`
	PushedAt			string		`json:"pushed_at"`
	GitLabUpdatedAt   	string 		`json:"last_activity_at,omitempty"`
	Language    		string 		`json:"language"`
	Languages			map[string]int	`json:"languages"`
	Archived    		bool   		`json:"archived"`
	OpenIssues  		int			`json:"open_issues_count"`
	Meta				Meta		`json:"meta"`
}

type Issue struct {
	ID        			uint64     	`json:"id"`
	Number    			uint64     	`json:"number"`
	GitLabNumber       	*uint64     `json:"iid,omitempty"`
	Description 		string		`json:"body"`
	GitlabDescription	string		`json:"description,omitempty"`
	Title     			string     	`json:"title"`
	User      			User       	`json:"user"`
	GitlabUser 			*GitlabUser `json:"author,omitempty"`
	State     			string     	`json:"state"`
	Assignees 			[]Assignee 	`json:"assignees"`
	Milestone 			*Milestone  `json:"milestone,omitempty"`
	Labels				[]Label     `json:"labels"`
	RepPath				string     `json:"reppath"`
	ProjectPath			string     `json:"project_path"`
	CreatedAt 			string     	`json:"created_at"`
	UpdatedAt 			string     	`json:"updated_at"`
	ClosedAt  			string     	`json:"closed_at"`
	HtmlUrl   			string     	`json:"html_url"`
	GitLabHTMLUrl     	string 		`json:"web_url,omitempty"`
	PullRequest			PullRequest `json:"pull_request"`
}

type User struct {
	ID           		uint64		`json:"id"`
	Login		 		string		`json:"login"`
	AvatarURL    		string 		`json:"avatar_url"`
	URL		 	 		string		`json:"html_url"`
}

type GitlabUser struct {
	ID           		uint64		`json:"id"`
	GitLabLogin		 	string		`json:"name"`
	AvatarURL    		string 		`json:"avatar_url"`
	GitLabHTMLUrl     	string 		`json:"web_url,omitempty"`
}

type Assignee struct {
	ID           		uint64		`json:"id"`
	Login		 		string		`json:"login"`
	AvatarURL    		string 		`json:"avatar_url"`
	URL	 	 	 		string		`json:"html_url"`
}

type Milestone struct {
	ID        			uint64     	`json:"id"`
	Number    			uint64     	`json:"number"`
	State     			string     	`json:"state"`
	Title     			string     	`json:"title"`
	Description 		string 		`json:"description"`
	Creator      		User    	`json:"creator"`
}

type Response struct {
	Repositories []Repos
	PageCount	int
}

type Label struct {
	ID          int    `json:"id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	URL         string `json:"url,omitempty"`
	Name        string `json:"name"`
	Type		string	`json:"type"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type PullRequest struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}