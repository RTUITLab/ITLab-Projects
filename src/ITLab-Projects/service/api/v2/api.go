package v2

import (
	"context"

	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/api/v2/projects"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type Api struct {
	Repository		*repositories.Repositories
	Testmode		bool
	Logger			log.Logger
	Auth			endpoint.Middleware

	projectService	projects.Service
}

type Config struct {
	Testmode	bool
	Config		config.AuthConfig
}

type ApiEndpoints struct {
	Projects	projects.Endpoints
}

func New(
	cfg 		Config,
	Repository	*repositories.Repositories,
) *Api {
	a := &Api{
		Repository: Repository,
	}
	a.Testmode = cfg.Testmode

	return a
}

func (a *Api) AddLogger(logger log.Logger) {
	a.Logger = logger
}

func (a *Api) AddAuthMiddleware(auth endpoint.Middleware) {
	a.Auth = auth
}

func (a *Api) CreateServices() {
	RepoImp := repoimpl.New(a.Repository)
	a.projectService = projects.New(
		RepoImp,
		a.Logger,
	)
}

func (a *Api) Build(r *mux.Router) {
	v2 := r.PathPrefix("/v2").Subrouter()

	var endpoints ApiEndpoints
	if a.Testmode {
		endpoints = a._buildEndpoints()
	} else {
		endpoints = a.buildEndpoints()
	}

	projects.NewHTTPServer(
		context.Background(),
		endpoints.Projects,
		v2,
	)
}