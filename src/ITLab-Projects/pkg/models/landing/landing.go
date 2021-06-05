package landing

import (
	"github.com/Kamva/mgm"
)

type LandingCompact struct {
	mgm.DefaultModel			`json:"-" bson:",inline" swaggerignore:"true"`
	RepoId			uint64		`json:"id" bson:"repo_id"`
	Title			string		`json:"title"`
	Image			[]string	`json:"images"`
	Date			Time		`json:"date" bson:"date"`
	Tags			[]string	`json:"tags"`
}

type Landing struct {
	LandingCompact				`bson:",inline"`
	Description	string			`json:"description"`
	Videos		[]string		`json:"videos"`
	Tech		[]string		`json:"tech"`
	Developers	[]string		`json:"developers"`
	Site		string			`json:"site"`
	SourceCode	[]*SourceCode	`json:"sourceCode"`
}

type SourceCode	struct {
	Name		string		`json:"name"`
	// repository link
	Value		string		`json:"link"`
}