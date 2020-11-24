package server

import (
	"ITLab-Projects/config"
	"ITLab-Projects/server/utils"
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Client
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
	log.WithField("dburi", cfg.DB.URI).Info("Current database URI: ")
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		log.WithFields(log.Fields{
			"function": "mongo.NewClient",
			"error":    err,
			"db_uri":   cfg.DB.URI,
		},
		).Fatal("Failed to create new MongoDB client")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "mongo.Connect",
			"error":    err},
		).Fatal("Failed to connect to MongoDB")
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "mongo.Ping",
			"error":    err},
		).Fatal("Failed to ping MongoDB")
	}
	log.Info("Connected to MongoDB!")

	dbName := utils.GetDbName(cfg.DB.URI)
	log.WithFields(log.Fields{
		"db_name": dbName,
	}).Info("Database information: ")
	log.WithField("testMode", cfg.App.TestMode).Info("Let's check if test mode is on...")

	projectsCollection = client.Database(dbName).Collection("projects")
	repsCollection = client.Database(dbName).Collection("repos")
	labelsCollection = client.Database(dbName).Collection("labels")
	issuesCollection = client.Database(dbName).Collection("issues")

	a.Router = mux.NewRouter().UseEncodedPath()
	a.setRouters()
}

func (a *App) setRouters() {
	// TODO calc hash from secret and payload
	github := a.Router.PathPrefix("/api/projects").Subrouter()
	github.HandleFunc("/update", updateInfo).Methods("POST")

	private := a.Router.PathPrefix("/api/projects").Subrouter()
	if cfg.App.TestMode {
		private.Use(loggingMiddleware)
	} else {
		private.Use(authMiddleware)
	}

	private.HandleFunc("/forceupdate", forceUpdateInfo).Methods("POST")
	private.HandleFunc("/projects", getFilteredProjects).Methods("GET").Queries("filter", "{filter}", "labels", "{labels}")
	private.HandleFunc("/projects", getFilteredProjects).Methods("GET").Queries("labels", "{labels}")
	private.HandleFunc("/projects", getFilteredProjects).Methods("GET").Queries("filter", "{filter}")
	private.HandleFunc("/projects", getAllProjects).Methods("GET")
	private.HandleFunc("/projects/{path}", getProjectReps).Methods("GET")
	private.HandleFunc("/labels", getAllLabels).Methods("GET")
	private.HandleFunc("/reps", getRepsPage).Methods("GET").Queries("page", "{page}")
	private.HandleFunc("/reps/{id}", getRep).Methods("GET").Queries("platform", "{platform}")
	private.HandleFunc("/reps/{id}/issues", getAllIssuesForRep).Methods("GET").Queries("platform", "{platform}", "state", "{state}")
	private.HandleFunc("/issues", getFilteredIssues).Methods("GET").Queries("filter", "{filter}", "labels", "{labels}")
	private.HandleFunc("/issues", getFilteredIssues).Methods("GET").Queries("labels", "{labels}")
	private.HandleFunc("/issues", getFilteredIssues).Methods("GET").Queries("filter", "{filter}")
	private.HandleFunc("/issues", getAllOpenedIssues).Methods("GET")
	private.HandleFunc("/issues/{reppath}", getProjectIssues).Methods("GET")
	private.HandleFunc("/reps/{id}/issues/{number}", getIssue).Methods("GET").Queries("platform", "{platform}")
}

func (a *App) Run(addr string) {
	log.WithField("port", cfg.App.AppPort).Info("Starting server on port:")
	log.Info("\n\nNow handling routes!")

	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "http.ListenAndServe",
			"error":    err},
		).Fatal("Failed to run a server!")
	}
}
