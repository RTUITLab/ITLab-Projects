package v1

import (
	"net/url"
	"strconv"
	"github.com/ITLab-Projects/pkg/apibuilder"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type beforeDeleteFunc func(interface{}) error

type Api struct {
	Repository *repositories.Repositories
	Requester githubreq.Requester
}

func New(
	Repository *repositories.Repositories,
	Requester githubreq.Requester,
	) apibuilder.ApiBulder {
	return &Api{
		Repository: Repository,
		Requester: Requester,
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
	// TODO
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