package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func getAllReps(w http.ResponseWriter, r *http.Request) {
	//reps := make([]models.Repos, 0)
	data := mux.Vars(r)
	tempReps, _ := getRepsFromGithub(data["page"])
	/*reps = append(reps, tempReps...)

	tempReps, gitlabPagesCount :=  getRepsFromGitlab(data["page"])
	reps = append(reps, tempReps...)

	if githubPagesCount > gitlabPagesCount {
		w.Header().Set("X-Total-Pages", strconv.Itoa(githubPagesCount))
	} else {
		w.Header().Set("X-Total-Pages", strconv.Itoa(gitlabPagesCount))
	}
*/
	json.NewEncoder(w).Encode(tempReps)
}

func getRep(w http.ResponseWriter, r *http.Request) {
	var rep models.Repos

	data := mux.Vars(r)
	platform := data["platform"]

	switch platform {
	case "github":
		rep = getRepFromGithub(data["id"])
	case "gitlab":
		rep = getRepFromGitlab(data["id"])
	}
	json.NewEncoder(w).Encode(rep)
}

func getAllIssues(w http.ResponseWriter, r *http.Request) {
	var issues []models.Issue
	data := mux.Vars(r)
	platform := data["platform"]

	switch platform {
	case "github":
		issues = getIssuesForGithub(data["id"], data["state"])
	case "gitlab":
		issues = getIssuesForGitlab(data["id"], data["state"])
	}
	json.NewEncoder(w).Encode(issues)
}

func getIssue(w http.ResponseWriter, r *http.Request) {
	var issue models.Issue

	data := mux.Vars(r)
	platform := data["platform"]

	switch platform {
	case "github":
		issue = getIssueFromGithub(data["id"], data["number"])
	case "gitlab":
		issue = getIssueFromGitlab(data["id"], data["number"])
	}

	json.NewEncoder(w).Encode(issue)
}

