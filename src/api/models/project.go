package models

type Project struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
	LastUpdated		string			`json:"lastUpdated"`
	Reps			[]string		`json:"reps"`
	StackTags		StackTags		`json:"stackTags"`
}

type ProjectInfo struct {
	Project		Project		`json:"project"`
	Repos		Meta		`json:"repos"`
}

type Meta struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
	StackTags		StackTags		`json:"stackTags"`
}

type StackTags struct {
	Directions		[]string		`json:"directions"`
	Databases		[]string		`json:"databases"`
	Frameworks		[]string		`json:"frameworks"`

}