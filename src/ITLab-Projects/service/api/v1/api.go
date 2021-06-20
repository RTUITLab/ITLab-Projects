package v1

import (
	updateS "github.com/ITLab-Projects/service/api/v1/updater"
	"context"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strconv"
	"time"

	"github.com/ITLab-Projects/service/api/v1/landing"

	"github.com/go-kit/kit/endpoint"
	kit_logger "github.com/go-kit/kit/log"
	"github.com/ITLab-Projects/service/api/v1/estimate"
	"github.com/ITLab-Projects/service/api/v1/functask"
	"github.com/ITLab-Projects/service/repoimpl"

	"github.com/ITLab-Projects/service/api/v1/issues"
	"github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/ITLab-Projects/service/api/v1/tags"

	"github.com/ITLab-Projects/pkg/updater"

	_ "github.com/ITLab-Projects/docs"
	"github.com/ITLab-Projects/pkg/config"
	swag "github.com/swaggo/http-swagger"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	Repository 		*repositories.Repositories
	Requester 		githubreq.Requester
	MFSRequester	mfsreq.Requester
	Testmode		bool
	upd				*updater.Updater
	NewAuth			endpoint.Middleware

	projectService	projects.Service
	issueService	issues.Service
	tagsService		tags.Service
	taskService		functask.Service
	estService		estimate.Service
	landingService	landing.Service
	updaterService	updateS.Service
	Logger			kit_logger.Logger
}

type Config struct {
	Testmode 		bool
	UpdateTime		string
	Config config.AuthConfig
}

type ApiEndpoints struct {
	Issues 		issues.Endpoints
	Projects 	projects.Endpoints
	Tags		tags.Endpoints
	Task		functask.Endpoints
	Est			estimate.Endpoints
	Landing		landing.Endpoints
	Update		updateS.Endpoints
}

func New(
	cfg Config,
	Repository *repositories.Repositories,
	Requester githubreq.Requester,
	MFSRequester	mfsreq.Requester,
) *Api {
	a := &Api{
		Repository: Repository,
		Requester: Requester,
		MFSRequester: MFSRequester,
	}

	a.Testmode = cfg.Testmode
	if cfg.UpdateTime != "" {
		log.Debug("WithUpdater")
		a.WithUpdater(cfg.UpdateTime)
	}
	log.Debug(a.upd)

	return a
}

func (a *Api) AddLogger(logger kit_logger.Logger) {
	a.Logger = logger
}

func (a *Api) AddAuthMiddleware(auth endpoint.Middleware) {
	a.NewAuth = auth
}

func (a *Api) CreateServices() {
	RepoImp := repoimpl.New(a.Repository)
	a.projectService = projects.New(
		RepoImp,
		a.Logger,
		a.MFSRequester,
	)

	a.estService = estimate.New(
		RepoImp,
		a.Logger,
		a.MFSRequester,
	)

	a.issueService = issues.New(
		RepoImp,
		a.Logger,
	)

	a.tagsService = tags.New(
		RepoImp,
		a.Logger,
	)

	a.taskService = functask.New(
		RepoImp,
		a.MFSRequester,
		a.Logger,
	)

	a.landingService = landing.New(
		RepoImp,
		a.Logger,
	)

	a.updaterService = updateS.New(
		RepoImp,
		a.Logger,
		a.Requester,
		a.upd,
	)
}

func (a *Api) WithUpdater(Time string) *Api {
	if err := a.CreateUpdater(Time); err != nil {
		log.WithFields(
			log.Fields{
				"package": "api/v1",
				"func": "WithUpdater",
				"err": err,
			},
		).Panic("Failed to create with updater")
	}

	return a
}

func (a *Api) CreateUpdater(Time string) error {
	duration, err := time.ParseDuration(Time)
	if err != nil {
		return err
	}

	a.upd = updater.New(
		duration,
		a.update,
	)

	return nil
}

func (a *Api) StartUpdater() {
	if a.upd != nil {
		go a.upd.Update()
	}
}

func (a *Api) resetUpdater() {
	if a.upd != nil {
		log.Debug("Reset update")
		a.upd.Reset()
	}
}

func (a *Api) update() {
	log.Debug("Start update")
	sessctx, err := repositories.GetMongoSessionContext(context.Background())
	if err != nil {
		log.WithFields(
			log.Fields{
				"package": "api/v1",
				"func": "update",
				"err": err,
			},
		).Error("Failed to update projects")
	}
	log.Debug("put session")
	ctx := updater.WithUpdateContext(sessctx)
	if err := a.updaterService.UpdateProjects(ctx); err != nil {
		log.WithFields(
			log.Fields{
				"package": "api/v1",
				"func": "update",
				"err": err,
			},
		).Error("Failed to update projects")
	}
	log.Debug("End update")
}


func (a *Api) Build(r *mux.Router) {
	docs := r.PathPrefix("/swagger")

	projectsR := r.PathPrefix("/v1").Subrouter()

	// Docs
	docs.Handler(
		swag.WrapHandler,
	)

	var endpoints ApiEndpoints
	if a.Testmode {
		endpoints = a._buildEndpoint()
	} else {
		endpoints = a.buildEndpoints()
	}

	projects.NewHTTPServer(
		context.Background(),
		endpoints.Projects,
		projectsR,
	)

	issues.NewHTTPServer(
		context.Background(),
		endpoints.Issues,
		projectsR,
	)

	tags.NewHTTPServer(
		context.Background(),
		endpoints.Tags,
		projectsR,
	)

	functask.NewHTTPServer(
		context.Background(),
		endpoints.Task,
		projectsR,
	)

	estimate.NewHTTPServer(
		context.Background(),
		endpoints.Est,
		projectsR,
	)

	landing.NewHTTPServer(
		context.Background(),
		endpoints.Landing,
		projectsR,
	)

	updateS.NewHTTPServer(
		context.Background(),
		endpoints.Update,
		projectsR,
	)

	if a.Testmode {
		r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
		
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	a.StartUpdater()
}



func getUint(v url.Values, name string) uint64 {
	_value := v.Get(name)

	if _value == "" {
		return 0
	}

	value, err := strconv.ParseUint(_value, 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func logError(message, Handler string, err error) {
	prepare(Handler, err).Error(message)
}

func prepare(Handler string, err error) *log.Entry {
	return log.WithFields(
		log.Fields{
			"package": "api/v1",
			"handler": Handler,
			"err": err,
		},
	)
}

func catchPanic() {
	if r := recover(); r != nil {
		log.WithFields(
			log.Fields{
				"package": "api/v1",
				"panic": r,
			},
		).Info()
	}
}