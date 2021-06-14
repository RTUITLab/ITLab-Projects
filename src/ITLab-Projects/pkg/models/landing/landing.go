package landing

import (
	"bytes"

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
	Site		Site			`json:"site"`
	SourceCode	[]*SourceCode	`json:"sourceCode"`
}

type SourceCode	struct {
	Name		string		`json:"name"`
	// repository link
	Value		string		`json:"link"`
}

type Site string

func (s Site) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}

	if len(string(s)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(s) + `"`)
	}
	return buf.Bytes(), nil
}

func (s *Site) UnmarshalJSON(data []byte) error {
	str := string(data)

	if str == `null` {
		*s = ""
		return nil
	}
	res := Site(str)

	if len(res) >= 2 {
		res = res[1:len(res)-1]
	}
	*s = res
	return nil
}