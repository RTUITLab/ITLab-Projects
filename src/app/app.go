package app

import (
	"fmt"
	"net/http"

	"github.com/ITLab-Projects/pkg/apibuilder"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/repositories"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type App struct {
	Router 	*mux.Router
	Repository *repositories.Repositories
	Requester githubreq.Requester
	Port 	string
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

	app.Router = mux.NewRouter()

	return app
}

func (a *App) AddApi(Builders ...apibuilder.ApiBulder) {
	for _, Builder := range Builders {
		Builder.Build(a.Router)
	}
}

func (a *App) Start() {
	log.Infof("Starting Application is port %s", a.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s",a.Port), a.Router); err != nil {
		log.Panicf("Failed to start application %v", err)
	}
}



