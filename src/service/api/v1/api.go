package v1

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	requester.HandleFunc("/task", a.AddFuncTask).Methods("POST")
	requester.HandleFunc("/estimate", a.AddEstimate).Methods("POST")
	requester.HandleFunc("/task/{milestone_id:[0-9]+}", a.DeleteFuncTask).Methods("DELETE")
	requester.HandleFunc("/estimate/{milestone_id:[0-9]+}", a.DeleteEstimate).Methods("DELETE")
}

// UpdateAllProjects
// 
// @Summary Update all projects
// 
// @Description make all request to github to update repositories, milestones
// 
// @Description If don't get from gh some repos delete it in db
// 
// @Router /api/v1/projects/ [post]
// 
// @Success 200
// 
// @Failure 502 {object} e.Err
// 
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
					Message: "Can't get repositories",
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
		msChan <- ms
	}()

	go func() {
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
		rsChan <- rs
	}()

	go func() {
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
		tgsChan <- tgs
	}()
	var (
		ms 	[]milestone.MilestoneInRepo = nil
		rs 	[]realese.RealeseInRepo		= nil
		tgs	[]tag.Tag					= nil
	)

	// TODO Переделать работу с каналами

	for i := 0; i < 3; i++ {
		select {
		case <- ctx.Done():
			err := <- errChan
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(
				e.Err{
					Err: err.Error(),
					Message: e.Message {
						Message: "Failed to update",
					},
				},
			)
			return
		case _ms := <- msChan:
			ms = _ms
		case _rs := <- rsChan:
			rs = _rs
		case _tgs := <- tgsChan:
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
	// TODO delete also functask and estimate for deleted milestones
}

// AddFuncTask
// 
// @Summary add func task to milestone
// 
// @Description add func task to milestone
// 
// @Description if func task is exist for milesotne will replace it
// 
// @Router /api/v1/projects/task [post]
// 
// @Accept json
// 
// @Produce json
// 
// @Param functask body functask.FuncTask true "function task that you want to add"
// 
// @Success 201
// 
// @Failure 400 {object} e.Message "Unexpected body"
// 
// @Failure 500 {object} e.Message "Failed to save functask"
// 
// @Failure 404 {object} e.Message "Don't find milestone with this id"
func (a *Api) AddFuncTask(w http.ResponseWriter, r *http.Request) {
	var fntask functask.FuncTask
	if err := json.NewDecoder(r.Body).Decode(&fntask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Unexpected body",
			},
		)
		return
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := a.Repository.Milestone.GetOne(
		ctx,
		bson.M{"id": fntask.MilestoneID},
		func(sr *mongo.SingleResult) error {
			return nil
		},
		options.FindOne(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Don't find milestone with this id",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to save funtask",
			},
		)
		prepare("AddFuncTask", err).Error()
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

	w.WriteHeader(http.StatusCreated)
}

// DeleteFuncTask
// 
// @Summary delete functask from database
// 
// @Description delete functask from database
// 
// @Router /api/v1/projects/task/{milestone_id} [delete]
// 
// @Param milestone_id path uint64 true "should be uint"
// 
// @Produce json
// 
// @Success 200
// 
// @Failure 404 {object} e.Message "func task not found"
// 
// @Failure 500 {object} e.Message "Failed to delete func task"
func (a *Api) DeleteFuncTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	ctx, _ := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)

	if err := a.Repository.FuncTask.DeleteOne(
		ctx,
		bson.M{"milestone_id": milestoneID},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}
			return nil
		},
		options.Delete(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "func task not found",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to delete functask",
			},
		)
		prepare("DeleteFuncTask", err).Error()
		return
	}
}

// AddEstimate
// 
// @Summary add estimate to milestone
// 
// @Description add estimate to milestone
// 
// @Description if estimate is exist for milesotne will replace it
// 
// @Router /api/v1/projects/estimate [post]
// 
// @Accept json
// 
// @Produce json
// 
// @Param estimate body estimate.Estimate true "estimate that you want to add"
// 
// @Success 201
// 
// @Failure 400 {object} e.Message "Unexpected body"
// 
// @Failure 500 {object} e.Message "Failed to save estimate"
// 
// @Failure 404 {object} e.Message "Don't find milestone with this id"
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

	ctx, _ := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)

	if err := a.Repository.Milestone.GetOne(
		ctx,
		bson.M{"id": est.MilestoneID},
		func(sr *mongo.SingleResult) error {
			return nil
		},
		options.FindOne(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Don't find milestone with this id",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to save estimate",
			},
		)
		prepare("AddEstimate", err).Error()
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

	w.WriteHeader(http.StatusCreated)
}

// DeleteEstimate
// 
// @Summary delete estimate from database
// 
// @Description delete estimate from database
// 
// @Router /api/v1/projects/estimate/{milestone_id} [delete]
// 
// @Param milestone_id path uint64 true "should be uint"
// 
// @Produce json
// 
// @Success 200
// 
// @Failure 404 {object} e.Message "estimate not found"
// 
// @Failure 500 {object} e.Message "Failed to delete estimate"
func (a *Api) DeleteEstimate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	ctx, _ := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)

	if err := a.Repository.Estimate.DeleteOne(
		ctx,
		bson.M{"milestone_id": milestoneID},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}

			return nil
		},
		options.Delete(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "estimate not found",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to delete estimate",
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