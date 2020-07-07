package models

type Project struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
	Reps			[]string		`json:"reps"`
}

type ProjectInfo struct {
	Project		Meta		`json:"project"`
	Repos		Meta		`json:"repos"`
}

type Meta struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
}

