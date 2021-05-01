package realese

type Realese struct {
	ID			uint64		`json:"id"`
	HTMLURL 	string		`json:"html_url"`
	URL			string		`json:"url"`
}

type RealeseInRepo struct {
	RepoID	uint64	`json:"repo_id"`
	Realese 		`bson:",inline"`
}