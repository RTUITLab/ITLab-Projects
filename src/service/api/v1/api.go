package v1

import (
	"github.com/ITLab-Projects/pkg/models/tag"
	"context"
	"encoding/json"
	"net/http"
	"time"

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
		logError("Can't repositories", "UpdateAllProjects", err)
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

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	go func() {
		ms := a.Requester.GetAllMilestonesForRepoWithID(
			repo.ToRepo(repos),
			func(e error) {
				if e == githubreq.ErrForbiden || e == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Failed to get milestone")
					cancel()
				}
				prepare("UpdateAllProjects", err).Warn("Failed to get milestone")
			},
		)
		msChan <- ms
		close(msChan)
	}()
	go func() {
		rs := a.Requester.GetLastsRealeseWithRepoID(
			repo.ToRepo(repos),
			func(e error) {
				if e == githubreq.ErrForbiden || e == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Failed to get realese")
					cancel()
				}
				prepare("UpdateAllProjects", err).Warn("Failed to get realese")
			},
		)
		rsChan <- rs
		close(rsChan)
	}()

	go func() {
		tgs := a.Requester.GetAllTagsForRepoWithID(
			repo.ToRepo(repos),
			func(e error) {
				if e == githubreq.ErrForbiden || e == githubreq.ErrUnatorizared {
					prepare("UpdateAllProjects", err).Error("Faield to get tag")
					cancel()
				}
				prepare("UpdateAllProjects", err).Warn("Faield to get tag")
			},
		)
		tgsChan <- tgs
		close(tgsChan)
	}()
	var (
		ms 	[]milestone.MilestoneInRepo = nil
		rs 	[]realese.RealeseInRepo		= nil
		tgs	[]tag.Tag					= nil
	)

	for {
		select {
		case _ms := <- msChan:
			ms = _ms
		case _rs := <- rsChan:
			rs = _rs
		case _tgs := <-tgsChan:
			tgs = _tgs
		case <- ctx.Done():
			return
		default:
			if ms != nil && rs != nil && tgs != nil {
				break
			}
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

// TODO add handler to add estimeate
// TODO add handler to add func_task
// TODO add handler to get projs by chunks
// TODO 

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