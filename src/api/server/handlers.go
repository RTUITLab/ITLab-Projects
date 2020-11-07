package server

import (
	"ITLab-Projects/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getProjectReps(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	var reps []models.Repos
	data := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"path" : data["path"]}
	err := projectsCollection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.FindOne",
			"handler" : "getProjectReps",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error, shutting down...")
		w.WriteHeader(404)
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter = bson.M{"path" : bson.M{"$in" : project.Reps}}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"path": 1})
	cur, err := repsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getProjectReps",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error, shutting down...")
		w.WriteHeader(404)
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reps)
	json.NewEncoder(w).Encode(reps)
}

func getRepsPage(w http.ResponseWriter, r *http.Request) {
	reps := make([]models.Repos, 0)
	data := mux.Vars(r)
	pageNum, err := strconv.Atoi(data["page"])
	if err != nil || pageNum < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repsTotal, err := repsCollection.CountDocuments(ctx, bson.M{})
	pageTotal := calcPageTotal(repsTotal)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	opts := options.Find()
	opts.SetLimit(int64(cfg.App.ElemsPerPage))
	opts.SetSort(bson.M{"path" : 1})
	opts.SetSkip(int64((pageNum-1) * cfg.App.ElemsPerPage))
	cur, err := repsCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getPageRepsFromGithub",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reps)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getPageRepsFromGithub",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}

	w.Header().Set("X-Total-Pages", strconv.Itoa(pageTotal))
	json.NewEncoder(w).Encode(reps)
}

func getFilteredProjects(w http.ResponseWriter, r *http.Request) {
	projects := make([]models.Project, 0)
	data := mux.Vars(r)
	filterTag := data["filter"]
	labelsFilter := strings.Split(data["labels"], "&")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	opts := options.Find()
	opts.SetSort(bson.M{"path" : 1})

	filter := bson.M{"path" : bson.M{"$regex" : primitive.Regex{Pattern: filterTag, Options: "i"}}}
	if data["labels"] != "" {
		filter = bson.M{"path" : bson.M{"$regex" : primitive.Regex{Pattern: filterTag, Options: "i"}},
			"labels.name" : bson.M{"$all" : labelsFilter}}
	}

	cur, err := projectsCollection.Find(ctx, filter, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getFilteredProjects",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &projects)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getFilteredProjects",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(projects)
}

func getFilteredIssues(w http.ResponseWriter, r *http.Request) {
	issues := make([]models.Issue, 0)
	data := mux.Vars(r)
	filterTag := data["filter"]
	labelsFilter := strings.Split(data["labels"], "&")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	opts := options.Find()
	opts.SetSort(bson.M{"title" : 1})

	filter := bson.M{"title" : bson.M{"$regex" : primitive.Regex{Pattern: filterTag, Options: "i"}},
		"pullrequest.url": "", "state": "open"}
	if data["labels"] != "" {
		log.Info(filterTag)
		log.Info(labelsFilter)
		filter = bson.M{"title" : bson.M{"$regex" : primitive.Regex{Pattern: filterTag, Options: "i"}},
			"labels.name" : bson.M{"$all" : labelsFilter}, "pullrequest.url": "", "state": "open"}
	}

	cur, err := issuesCollection.Find(ctx, filter, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getFilteredIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &issues)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getFilteredIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(issues)
}

func getAllOpenedIssues(w http.ResponseWriter, r *http.Request) {
	issues := make([]models.Issue, 0)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"updatedat": -1})
	filter := bson.M{"pullrequest.url": "", "state": "open"}
	cur, err := issuesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &issues)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getAllIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(issues)
}

func getProjectIssues(w http.ResponseWriter, r *http.Request) {
	issues := make([]models.Issue, 0)
	data := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"updatedat": -1})
	filter := bson.M{"$and" : []bson.M{
		{"pullrequest.url": ""},
		{"projectpath": data["reppath"]},
	}}
	cur, err := issuesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getProjectIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &issues)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getProjectIssues",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(issues)
}
func getRep(w http.ResponseWriter, r *http.Request) {
	var rep models.Repos
	data := mux.Vars(r)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"path" : data["id"]}
	err := repsCollection.FindOne(ctx, filter).Decode(&rep)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.FindOne",
			"handler" : "getRep",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	fmt.Println(rep.Meta.Description)
	json.NewEncoder(w).Encode(rep)
}

func getAllIssuesForRep(w http.ResponseWriter, r *http.Request) {
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

func getAllProjects(w http.ResponseWriter, r *http.Request) {
	projects := make([]models.Project, 0)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"lastUpdated": -1})
	cur, err := projectsCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllProjects",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &projects)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getAllProjects",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(projects)
}

func forceUpdateInfo(w http.ResponseWriter, r *http.Request) {
	cGithub := make(chan models.Response)
	cProjects := make(chan models.ProjectInfo)
	var projects []models.ProjectInfo

	go getRepsFromGithub("all", cGithub)
	result := <-cGithub
	for i, _ := range result.Repositories {
		go getProjectInfoFile(&result.Repositories[i], cProjects)
	}
	for i := 0; i< len(result.Repositories); i++  {
		project := <-cProjects
		if project.Project.Path != "" {
			projects = append(projects, project)
		}
	}
	saveReposToDB(result.Repositories)
	saveLabelsToDB(result.Repositories)
	log.Info("Performed data update!!! ")
	w.WriteHeader(200)
}

func getAllLabels(w http.ResponseWriter, r *http.Request) {
	labels := make([]models.Label, 0)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"type": 1})
	cur, err := labelsCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllLabels",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &labels)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getAllLabels",
			"error"	:	err,
		},
		).Warn("DB interaction resulted in error")
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(labels)
}

func updateInfo(w http.ResponseWriter, r *http.Request) {
	var payload	models.WebhookPayload
	json.NewDecoder(r.Body).Decode(&payload)
	switch {
	case payload.Issue.ID != 0:
		saveToDB(payload.Issue)
	case payload.Label.ID != 0:
		saveLabelToDB(payload.Label)
	case payload.Repository.ID != 0 || payload.Ref != "":
		saveToDB(payload.Repository)
	default:
		w.WriteHeader(404)
	}
	w.WriteHeader(200)
}