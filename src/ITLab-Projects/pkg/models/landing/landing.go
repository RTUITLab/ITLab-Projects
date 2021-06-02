package landing

import (
	"time"

	"github.com/Kamva/mgm"
)

type LandingCompact struct {
	mgm.DefaultModel			`json:"-" bson:",inline" swaggerignore:"true"`
	RepoId			uint64		`json:"repo_id" bson:"repo_id"`
	Title			string		`json:"title"`
	Image			[]string	`json:"image"`
	Date			time.Time	`json:"date"`
	Tags			[]string	`json:"tags"`
}

type Landing struct {
	LandingCompact				`bson:",inline"`
	Videos		[]string		`json:"videos"`
	Tech		[]string		`json:"tech"`
	Developers	[]string		`json:"developers"`
	Site		string			`json:"site"`
	SourceCode	[]*SourceCode	`json:"source_code"`
}

type SourceCode	struct {
	Name		string		`json:"name"`
	// repository link
	Value		string		`json:"value"`
}