package v1

import (
	"github.com/ITLab-Projects/service/middleware/contenttype"
	"github.com/ITLab-Projects/pkg/updater"
	"context"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strconv"
	"time"

	_ "github.com/ITLab-Projects/docs"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	swag "github.com/swaggo/http-swagger"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type beforeDeleteFunc func(
	interface{},
) error


type Api struct {
	Repository 		*repositories.Repositories
	Requester 		githubreq.Requester
	MFSRequester	mfsreq.Requester
	Testmode		bool
	upd				*updater.Updater
	Auth auth.AuthMiddleware
}

type Config struct {
	Testmode 		bool
	UpdateTime		string
	Config config.AuthConfig
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

	a.Auth = auth.New(cfg.Config)
	a.Testmode = cfg.Testmode

	if cfg.UpdateTime != "" {
		log.Debug("WithUpdater")
		a = a.WithUpdater(cfg.UpdateTime)
	}

	return a
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
	if err := a.updateAllProjects(sessctx); err != nil {
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

	projects := r.PathPrefix("/api/v1/projects").Subrouter()
	admin := projects.NewRoute().Subrouter()

	admin.HandleFunc("/", a.UpdateAllProjects).Methods("POST")
	admin.HandleFunc("/task", a.AddFuncTask).Methods("POST")
	admin.HandleFunc("/estimate", a.AddEstimate).Methods("POST")
	admin.HandleFunc("/task/{milestone_id:[0-9]+}", a.DeleteFuncTask).Methods("DELETE")
	admin.HandleFunc("/estimate/{milestone_id:[0-9]+}", a.DeleteEstimate).Methods("DELETE")
	admin.HandleFunc("/{id:[0-9]+}", a.DeleteProject).Methods("DELETE")
	
	projects.HandleFunc("/", a.GetProjects).Methods("GET")
	projects.HandleFunc("/{id:[0-9]+}", a.GetProject).Methods("GET")
	projects.HandleFunc("/tags", a.GetTags).Methods("GET")
	projects.HandleFunc("/issues", a.GetIssues).Methods("GET")
	projects.HandleFunc("/issues/labels", a.GetLabels).Methods("GET")


	projects.Use(contenttype.AppJSON)
	if !a.Testmode {
		projects.Use(
			mux.MiddlewareFunc(a.Auth),
		)
		
		admin.Use(
			auth.AdminMiddleware,
		)

		docs.Handler(
			a.Auth(
				swag.WrapHandler,
			),
		)
	} else {
		docs.Handler(
			swag.WrapHandler,
		)
		r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
		
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	projects.Use(
		mgsess.PutSessionINTOCtx,
	)

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

func (a *Api) beforeDeleteWithReq(r *http.Request) beforeDeleteFunc {
	return func(v interface{}) error {
		return a.beforeDelete(
			a.MFSRequester.NewRequests(r),
			v,
		)
	}
}

func (a *Api) beforeDelete(
	deleter mfsreq.FileDeleter,
	v interface{},
	) error {
	log.Info("Before delete!")
	switch v.(type) {
	case estimate.EstimateFile:
		est, _ := v.(estimate.EstimateFile)
		if err := deleter.DeleteFile(est.FileID); err != nil {
			return err
		}
	case []estimate.EstimateFile:
		ests, _ := v.([]estimate.EstimateFile)
		for _, est := range ests {
			if err := deleter.DeleteFile(est.FileID); err != nil {
				return err
			}
		}
	case functask.FuncTaskFile:
		task, _ := v.(functask.FuncTaskFile)
		if err := deleter.DeleteFile(task.FileID); err != nil {
			return err
		}
	case []functask.FuncTaskFile:
		tasks, _ := v.([]functask.FuncTaskFile)
		for _, task := range tasks {
			if err := deleter.DeleteFile(task.FileID); err != nil {
				return err
			}
		}
	default:
	}
	return nil
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