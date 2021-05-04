package realese

import "github.com/Kamva/mgm"

type Realese struct {
	ID			uint64		`json:"id"`
	HTMLURL 	string		`json:"html_url"`
	URL			string		`json:"url"`
}

type RealeseInRepo struct {
	mgm.DefaultModel		`json:"-" bson:",inline"`
	RepoID	uint64			`json:"repo_id"`
	Realese 				`bson:",inline"`
}

func (r *RealeseInRepo) CollectionName() string {
	return "realese"
}