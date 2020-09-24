package server

import (
	"ITLab-Projects/config"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB *mongo.Client
}

var projectsCollection *mongo.Collection
var repsCollection *mongo.Collection
var labelsCollection *mongo.Collection
var issuesCollection *mongo.Collection
var cfg *config.Config
var httpClient *http.Client

func (a *App) Init(config *config.Config) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	cfg = config
	httpClient = createHTTPClient()
	log.Info("ITLab-Projects is starting up!")
	DBUri := "mongodb://" + cfg.DB.Host + ":" + cfg.DB.DBPort
	log.WithField("dburi", DBUri).Info("Current database URI: ")
	client, err := mongo.NewClient(options.Client().ApplyURI(DBUri))
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.NewClient",
			"error"	:	err,
			"db_uri":	DBUri,
		},
		).Warn("Failed to create new MongoDB client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Connect",
			"error"	:	err},
		).Warn("Failed to connect to MongoDB")
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Ping",
			"error"	:	err},
		).Warn("Failed to ping MongoDB")
	}
	log.Info("Connected to MongoDB!")
	log.WithFields(log.Fields{
		"db_name" : cfg.DB.DBName,
	}).Info("Database information: ")
	log.WithField("testMode", cfg.App.TestMode).Info("Let's check if test mode is on...")

	projectsCollection = client.Database(cfg.DB.DBName).Collection(cfg.DB.ProjectsCollectionName)
	repsCollection = client.Database(cfg.DB.DBName).Collection(cfg.DB.ReposCollectionName)
	labelsCollection = client.Database(cfg.DB.DBName).Collection(cfg.DB.LabelsCollectionName)
	issuesCollection = client.Database(cfg.DB.DBName).Collection(cfg.DB.IssuesCollectionName)

	a.Router = mux.NewRouter().UseEncodedPath()
	a.setRouters()
}

func (a *App) setRouters() {
	if cfg.App.TestMode {
		a.Router.Use(loggingMiddleware)
	} else {
		a.Router.Use(authMiddleware)
	}

	a.Router.HandleFunc("/api/projects/update", getRelevantInfo).Methods("POST")
	a.Router.HandleFunc("/api/projects/projects", getAllProjects).Methods("GET")
	a.Router.HandleFunc("/api/projects/projects/{path}", getProjectReps).Methods("GET")
	a.Router.HandleFunc("/api/projects/labels", getAllLabels).Methods("GET")
	a.Router.HandleFunc("/api/projects/reps", getFilteredReps).Methods("GET").Queries("filter","{filter}")
	a.Router.HandleFunc("/api/projects/reps", getRepsPage).Methods("GET").Queries("page","{page}")
	a.Router.HandleFunc("/api/projects/reps/{id}", getRep).Methods("GET").Queries("platform", "{platform}")
	a.Router.HandleFunc("/api/projects/reps/{id}/issues", getAllIssuesForRep).Methods("GET").Queries("platform", "{platform}", "state", "{state}")
	a.Router.HandleFunc("/api/projects/issues", getFilteredIssues).Methods("GET").Queries("filter","{filter}")
	a.Router.HandleFunc("/api/projects/issues", getAllOpenedIssues).Methods("GET")
	a.Router.HandleFunc("/api/projects/issues/{reppath}", getProjectIssues).Methods("GET")
	a.Router.HandleFunc("/api/projects/reps/{id}/issues/{number}", getIssue).Methods("GET").Queries("platform", "{platform}")
}

func (a *App) Run(addr string) {
	log.WithField("port", cfg.App.AppPort).Info("Starting server on port:")
	log.Info("\n\nNow handling routes!")

	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.ListenAndServe",
			"error"	:	err},
		).Fatal("Failed to run a server!")
	}
}