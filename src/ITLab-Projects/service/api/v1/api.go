package v1

import (
	"net/http"
	"net/http/pprof"
	"net/url"
	"strconv"

	_ "github.com/ITLab-Projects/docs"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	swag "github.com/swaggo/http-swagger"

	"github.com/ITLab-Projects/pkg/apibuilder"
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
	Auth auth.AuthMiddleware
}

type Config struct {
	Testmode bool
	Config config.AuthConfig
}

func New(
	cfg Config,
	Repository *repositories.Repositories,
	Requester githubreq.Requester,
	MFSRequester	mfsreq.Requester,
	) apibuilder.ApiBulder {
	a := &Api{
		Repository: Repository,
		Requester: Requester,
		MFSRequester: MFSRequester,
	}

	a.Auth = auth.New(cfg.Config)
	a.Testmode = cfg.Testmode

	return a
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