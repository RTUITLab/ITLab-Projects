package v1

import (
	"net/url"
	"strconv"

	"github.com/ITLab-Projects/pkg/apibuilder"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type beforeDeleteFunc func(interface{}) error

type Api struct {
	Repository 		*repositories.Repositories
	Requester 		githubreq.Requester
	MFSRequester	mfsreq.Requester
}

func New(
	Repository *repositories.Repositories,
	Requester githubreq.Requester,
	MFSRequester	mfsreq.Requester,
	) apibuilder.ApiBulder {
	return &Api{
		Repository: Repository,
		Requester: Requester,
		MFSRequester: MFSRequester,
	}
}

func (a *Api) Build(r *mux.Router) {
	requester := r.PathPrefix("/api/v1/projects").Subrouter()

	requester.HandleFunc("/", a.UpdateAllProjects).Methods("POST")
	requester.HandleFunc("/task", a.AddFuncTask).Methods("POST")
	requester.HandleFunc("/estimate", a.AddEstimate).Methods("POST")
	requester.HandleFunc("/task/{milestone_id:[0-9]+}", a.DeleteFuncTask).Methods("DELETE")
	requester.HandleFunc("/estimate/{milestone_id:[0-9]+}", a.DeleteEstimate).Methods("DELETE")
	requester.HandleFunc("/", a.GetProjects).Methods("GET")
	requester.HandleFunc("/{id:[0-9]+}", a.GetProject).Methods("GET")
	requester.HandleFunc("/{id:[0-9]+}", a.DeleteProject).Methods("DELETE")
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

func (a *Api) beforeDelete(v interface{}) error {
	log.Info("Before delete!")
	switch v.(type) {
	case estimate.EstimateFile:
		est, _ := v.(estimate.EstimateFile)
		if err := a.MFSRequester.DeleteFile(est.FileID); err != nil {
			return err
		}
	case []estimate.EstimateFile:
		ests, _ := v.([]estimate.EstimateFile)
		for _, est := range ests {
			if err := a.MFSRequester.DeleteFile(est.FileID); err != nil {
				return err
			}
		}
	case functask.FuncTaskFile:
		task, _ := v.(functask.FuncTaskFile)
		if err := a.MFSRequester.DeleteFile(task.FileID); err != nil {
			return err
		}
	case []functask.FuncTaskFile:
		tasks, _ := v.([]functask.FuncTaskFile)
		for _, task := range tasks {
			if err := a.MFSRequester.DeleteFile(task.FileID); err != nil {
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