package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getAllReps(w http.ResponseWriter, r *http.Request) {
	reps := make([]models.Repos, 0)
	w.Header().Set("Content-Type", "application/json")
	resp, err := http.Get("https://api.github.com/orgs/RTUITLab/repos")
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getAllReps",
			"error"	:	err,
		},
		).Warn("Can't reach GitHub API!")
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&reps)
	json.NewEncoder(w).Encode(reps)
}

func getRep(w http.ResponseWriter, r *http.Request) {
	var rep models.Repos
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	resp, err := http.Get("https://api.github.com/repos/RTUITLab/"+data["name"])
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRep",
			"error"	:	err,
		},
		).Warn("Can't reach GitHub API!")
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&rep)
	json.NewEncoder(w).Encode(rep)
}

func getAllIssues(w http.ResponseWriter, r *http.Request) {
	issues := make([]models.Issue, 0)
	var url string

	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	switch data["state"] {
	case "opened":
		url = "https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues?state=opened"
	case "closed":
		url = "https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues?state=closed"
	case "all":
		url = "https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues?state=all"
	default:
		url = "https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues"
	}
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getAllIssues",
			"error"	:	err,
		},
		).Warn("Can't reach GitHub API!")
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&issues)
	json.NewEncoder(w).Encode(issues)
}

func getOpenIssues(w http.ResponseWriter, r *http.Request) {
	issues := make([]models.Issue, 0)
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	resp, err := http.Get("https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues")
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getAllIssues",
			"error"	:	err,
		},
		).Warn("Can't reach GitHub API!")
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&issues)
	json.NewEncoder(w).Encode(issues)
}

func getIssue(w http.ResponseWriter, r *http.Request) {
	var issue models.Issue
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	resp, err := http.Get("https://api.github.com/repos/RTUITLab/"+data["name"]+"/issues/"+data["number"])
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getIssue",
			"error"	:	err,
		},
		).Warn("Can't reach GitHub API!")
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&issue)
	json.NewEncoder(w).Encode(issue)
}