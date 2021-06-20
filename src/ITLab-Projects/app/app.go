package app

import (
	klogrus "github.com/go-kit/kit/log/logrus"
	kl "github.com/go-kit/kit/log"
	"fmt"
	"net/http"
	"time"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/go-kit/kit/endpoint"

	"github.com/ITLab-Projects/pkg/apibuilder"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/repositories"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type App struct {
	Router 			*mux.Router
	Repository 		*repositories.Repositories
	Requester 		githubreq.Requester
	MFSRequester	mfsreq.Requester
	Port 			string
	Logger			kl.Logger

	auth			endpoint.Middleware
}

func New(cfg *config.Config) *App {
	app := &App{}

	app.Port = cfg.App.AppPort

	if _rep, err := repositories.New(&repositories.Config{
		DBURI: cfg.DB.URI,
	}); err != nil {
		log.WithFields(
			log.Fields{
				"package": "app",
				"func": "New",
				"err": err,
			},
		).Panic("Failed to init App")
	} else {
		app.Repository = _rep
	}

	app.Requester = githubreq.New(&githubreq.Config{
		AccessToken: cfg.Auth.Github.AccessToken,
	})

	app.MFSRequester = mfsreq.New(
		&mfsreq.Config{
			BaseURL: cfg.Services.MFS.BaseURL,
			TestMode: cfg.App.TestMode,
		},
	)

	app.Router = mux.NewRouter().PathPrefix("/api/projects").Subrouter()
	if !cfg.App.TestMode {
		log.SetLevel(log.InfoLevel)
		app.auth = auth.NewGoKitAuth(cfg.Auth)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	app.Logger = klogrus.NewLogrusLogger(log.StandardLogger())

	return app
}

func (a *App) AddApi(Builders ...apibuilder.ApiBulder) {
	for _, Builder := range Builders {
		Builder.AddLogger(a.Logger)
		Builder.AddAuthMiddleware(a.auth)
		Builder.CreateServices()
		Builder.Build(a.Router)
	}
}

func (a *App) Start() {
	log.Infof("Starting Application is port %s", a.Port)
	s := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Port),
		Handler: a.Router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout: 2*time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Panicf("Failed to start application %v", err)
	}
}



