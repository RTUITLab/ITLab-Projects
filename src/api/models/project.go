package models

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
func NewProjectInfo() ProjectInfo {
	projectInfo := ProjectInfo{}
	projectInfo.Repos.Labels = make([]Label, 0)
	projectInfo.Project.Labels = make([]Label, 0)
	return projectInfo
}