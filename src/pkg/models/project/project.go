package project

import (
	. "github.com/ITLab-Projects/pkg/models/label"
)

// TODO think aboyt need this or not

type Project struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
	LastUpdated		string			`json:"lastUpdated"`
	Reps			[]string		`json:"reps"`
	Labels			[]Label			`json:"labels"`
}

type ProjectInfo struct {
	Project		Project		`json:"project"`
	Repos		Meta		`json:"repos"`
}

type Meta struct {
	Path			string			`json:"path"`
	HumanName		string			`json:"humanName"`
	Description 	string			`json:"description"`
	Labels			[]Label			`json:"labels"`
}