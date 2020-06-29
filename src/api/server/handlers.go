package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getAllReps(w http.ResponseWriter, r *http.Request) {
	pageCount := 0
	reps := make([]models.Repos, 0)
	c := make(chan models.Response)
	result := make([]models.Response, 2)

	data := mux.Vars(r)
	go getRepsFromGithub(data["page"], c)
	go getRepsFromGitlab(data["page"], c)

	for i, _ := range result {
		result[i] = <-c
		reps = append(reps, result[i].Repositories...)
		if result[i].PageCount > pageCount {
			pageCount = result[i].PageCount
		}
	}
	w.Header().Set("X-Total-Pages", strconv.Itoa(pageCount))
	json.NewEncoder(w).Encode(reps)
}

func getAllRepsFromGithub(w http.ResponseWriter, r *http.Request) {
	c := make(chan models.Response)
	data := mux.Vars(r)
	go getRepsFromGithub(data["page"], c)
	result := <-c
	pageCount := result.PageCount
	w.Header().Set("X-Total-Pages", strconv.Itoa(pageCount))
	json.NewEncoder(w).Encode(result.Repositories)
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