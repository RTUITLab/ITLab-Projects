package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/pkg/errors"

	"github.com/ITLab-Projects/pkg/apibuilder"
	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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
	requester.HandleFunc("/add/functask", a.AddFuncTask).Methods("POST")
	requester.HandleFunc("/add/estimate", a.AddEstimate).Methods("POST")
	requester.HandleFunc("/delete/functask/{milestone_id:[0-9]+}", a.DeleteFuncTask).Methods("DELETE")
	requester.HandleFunc("/delete/estimate/{milestone_id:[0-9]+}", a.DeleteEstimate).Methods("DELETE")
}

// UpdateAllProjects
// @Summary Update all projects
// @Description make all request to github to update repositories, milestones
// @Description If don't get from gh some repos delete it in db
// @Router /api/v1/projects/ [post]
// @Success 200
// @Failure 502 {object} e.Err
// @Failure 500 {object} e.Message
func (a *Api) UpdateAllProjects(w http.ResponseWriter, r *http.Request) {
	repos, err := a.Requester.GetRepositories()
	if err == githubreq.ErrGetLastPage {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(
			e.Err{
				Err: err.Error(), 
				Message: e.Message {
					Message: "Try later we can't get pages of repo from githun",
				},
			},
		)
		logError("Can't get last page", "UpdateAllProjects", err)
		return
	} else if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
		logError("Can't get repositories", "UpdateAllProjects", err)
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(
			e.Err {
				Err: err.Error(),
				Message: e.Message{
					"Can't get repositories",
				},
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Can't get last page",
			},
		)
		logError("Can't get last page", "UpdateAllProjects", err)
		return
	}

	msChan := make(chan []milestone.MilestoneInRepo, 1)
	rsChan := make(chan []realese.RealeseInRepo, 1)
	tgsChan := make(chan []tag.Tag, 1)
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	go func() {
		log.Info("start ms")
		time.Sleep(1*time.Second)
		ms := a.Requester.GetAllMilestonesForRepoWithID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Failed to get milestone")
					cancel()
					errChan <- err
				} else if err != nil {
					prepare("UpdateAllProjects", err).Warn("Failed to get milestone")
				}
			},
		)
		log.Info("Send ms")
		msChan <- ms
	}()

	go func() {
		log.Info("start rs")
		time.Sleep(2*time.Second)
		rs := a.Requester.GetLastsRealeseWithRepoID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Failed to get realese")
					cancel()
					errChan <- err
				} else if errors.Unwrap(err) == githubreq.UnexpectedCode {

				} else if err != nil {
					prepare("UpdateAllProjects", err).Warn("Failed to get realese")
				}
			},
		)
		log.Info("Send rs")
		rsChan <- rs
	}()

	go func() {
		log.Info("start tgs")
		time.Sleep(3*time.Second)
		tgs := a.Requester.GetAllTagsForRepoWithID(
			ctx,
			repo.ToRepo(repos),
			func(err error) {
				if err == githubreq.ErrForbiden || err == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Faield to get tag")
					cancel()
					errChan <- err
				}
				prepare("UpdateAllProjects", err).Warn("Faield to get tag")
			},
		)
		log.Info("Send tgs")
		tgsChan <- tgs
	}()
	var (
		ms 	[]milestone.MilestoneInRepo = nil
		rs 	[]realese.RealeseInRepo		= nil
		tgs	[]tag.Tag					= nil
	)

	// TODO Переделать работу с каналами

	for i := 0; i < 3; i++ {
		log.Info("Start select")
		select {
		case <- ctx.Done():
			log.Info("Get cancel")
			err := <- errChan
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(
				e.Err{
					Err: err.Error(),
					Message: e.Message {
						"Failed to update",
					},
				},
			)
			return
		case _ms := <- msChan:
			log.Info("catch ms")
			ms = _ms
		case _rs := <- rsChan:
			log.Info("catch rs")
			rs = _rs
		case _tgs := <- tgsChan:
			log.Info("catch tgs")
			tgs = _tgs
		}
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := a.Repository.Repo.SaveAndDeletedUnfind(
		ctx,
		repo.ToRepo(repos),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e.Message{
			Message: "Can't save repositories",
		})
		prepare("UpdateAllProjects", err).Error("Can't save repositories")
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := a.Repository.Milestone.SaveAndDeletedUnfind(
		ctx,
		ms,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Can't save milestones",
			},
		)
		prepare("UpdateAllProjects", err).Error("Can't save milestones")
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := a.Repository.Realese.SaveAndDeletedUnfind(
		ctx,
		rs,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Can't save realeses",
			},
		)
		prepare("UpdateAllProjects", err).Error("Can't save realeses")
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := a.Repository.Tag.SaveAndDeletedUnfind(
		ctx,
		tgs,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Can't save realeses",
			},
		)
		prepare("UpdateAllProjects", err).Error("Can't save tags")
		return
	}
}

func (a *Api) AddFuncTask(w http.ResponseWriter, r *http.Request) {
	var fntask functask.FuncTask
	if err := json.NewDecoder(r.Body).Decode(&fntask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Unexpected body",
			},
		)
		prepare("AddFuncTask", err).Warn()
		return
	}

	if err := a.Repository.FuncTask.Save(fntask); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to save funtask",
			},
		)
		prepare("AddFuncTask", err).Error()
		return
	}
}

func (a *Api) DeleteFuncTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	if err := a.Repository.FuncTask.Delete(
		uint64(milestoneID),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				"Failed to delete functask",
			},
		)
		prepare("DeleteFuncTask", err).Error()
		return
	}
}

func (a *Api) AddEstimate(w http.ResponseWriter, r *http.Request) {
	var est estimate.Estimate
	if err := json.NewDecoder(r.Body).Decode(&est); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Unexpected body",
			},
		)
		prepare("AddEstimate", err).Warn()
		return
	}

	if err := a.Repository.Estimate.Save(est); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to save estimate",
			},
		)
		prepare("AddEstimate", err).Error()
		return
	}
}

func (a *Api) DeleteEstimate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	if err := a.Repository.Estimate.Delete(
		uint64(milestoneID),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				"Failed to delete estimate",
			},
		)
		prepare("DeleteEstimate", err).Error()
		return
	}
}


// TODO add handler to get projs by chunks

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