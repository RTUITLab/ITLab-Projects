package repo

import (
	. "github.com/ITLab-Projects/pkg/models/user"
)

type Repo struct {
	ID          		uint64 		`json:"id"`
	Name        		string 		`json:"name"`
	Contributors		[]User		`json:"contributors"`
	// Path				string		`json:"path_with_namespace,omitempty"`
	HTMLUrl     		string 		`json:"html_url"`
	Description 		string 		`json:"description"`
	CreatedAt   		string 		`json:"created_at"`
	UpdatedAt   		string 		`json:"updated_at"`
	PushedAt			string		`json:"pushed_at"`
	Language    		string 		`json:"language"`
	Languages			map[string]int	`json:"languages"`
	Archived    		bool   		`json:"archived"`
}

type RepoWithURLS struct {
	Repo
	LangaugesURL 		string 	`json:"languages_url"`
	ContributorsURL		string	`json:"contributors_url"`
}

func ToRepo(repos []RepoWithURLS) []Repo {
	var reps []Repo
	for _, rep := range repos {
		reps = append(reps, rep.Repo)
	}
	
	return reps
}