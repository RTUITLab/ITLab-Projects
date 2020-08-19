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
	"time"
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		).Fatal("DB interaction resulted in error, shutting down...")
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
		).Fatal("DB interaction resulted in error, shutting down...")
	}

	w.Header().Set("X-Total-Pages", strconv.Itoa(pageTotal))
	json.NewEncoder(w).Encode(reps)
}

func getFilteredReps(w http.ResponseWriter, r *http.Request) {
	reps := make([]models.Repos, 0)
	data := mux.Vars(r)
	filterTag := data["filter"]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	opts := options.Find()
	opts.SetSort(bson.M{"path" : 1})
	filter := bson.M{"path" : bson.M{"$regex" : primitive.Regex{Pattern: filterTag, Options: "i"}}}
	cur, err := repsCollection.Find(ctx, filter, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getPageRepsFromGithub",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
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
		).Fatal("DB interaction resulted in error, shutting down...")
	}

	json.NewEncoder(w).Encode(reps)
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
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	fmt.Println(rep.Meta.Description)
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

func getAllProjects(w http.ResponseWriter, r *http.Request) {
	projects := make([]models.Project, 0)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := projectsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllProjects",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
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
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	json.NewEncoder(w).Encode(projects)
}

func getRelevantInfo(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(200)
}
