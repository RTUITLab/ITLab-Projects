package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ITLab-Projects/pkg/models/repoasproj"

	"go.mongodb.org/mongo-driver/bson"

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
	requester.HandleFunc("/", a.GetProjects).Methods("GET")
	requester.HandleFunc("/{id:[0-9]+}", a.GetProject).Methods("GET")
}

// UpdateAllProjects
// 
// @Tags v1
// 
// @Summary Update all projects
// 
// @Description make all request to github to update repositories, milestones
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

	if err := a.Repository.Repo.Save(
		repo.ToRepo(repos),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e.Message{
			Message: "Can't save repositories",
		})
		prepare("UpdateAllProjects", err).Error("Can't save repositories")
		return
	}

	if err := a.Repository.Milestone.Save(
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

	if err := a.Repository.Realese.Save(
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

	if err := a.Repository.Tag.Save(
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
	
	// if err := a.Repository.FuncTask.DeleteFuncTasksNotIn(ms); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(
	// 		e.Message{
	// 			Message: "Can't delete unused func tasks",
	// 		},
	// 	)
	// 	prepare("UpdateAllProjects", err).Error("Can't delete unused func tasks")
	// 	return
	// }

	// if err := a.Repository.Estimate.DeleteEstimatesNotIn(ms); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(
	// 		e.Message{
	// 			Message: "Can't delete unused estimates",
	// 		},
	// 	)
	// 	prepare("UpdateAllProjects", err).Error("Can't delete unused estimates")
	// 	return
	// }
}

// AddFuncTask
// 
// @Tags v1 functask
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

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
// @Tags v1 functask
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)
	defer cancel()

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
// @Tags v1 estimate
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)
	defer cancel()

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
// @Tags v1 estimate
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)
	defer cancel()

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

// GetProjects
// 
// @Tags v1 projects
// 
// @Summary return projects according to query value
// 
// @Description return a projects you can filter count of them
// 
// @Description tags, name
// 
// @Produce json
// 
// @Router /api/v1/projects/ [get]
// 
// @Param start query integer false "represents the number of skiped projects"
// 
// @Param count query integer false "represent a limit of projects"
// 
// @Param tag query string false "use to filter projects by tag"
// 
// @Param name query string false "use to filter by name"
// 
// @Success 200 {array} repoasproj.RepoAsProjCompact
// 
// @Failure 500 {object} e.Message "Failed to get projects"
// 
// @Failure 500 {object} e.Message "Failed to get repositories"
func (a *Api) GetProjects(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	start := getUint(values, "start")
	count := getUint(values, "count")
	tag := values.Get("tag")
	name := values.Get("name")

	ctx := context.Background()

	if count == 0 {
		count = uint64(a.Repository.Repo.Count())
	}

	var repos []repo.Repo
	filter := bson.M{}
	if tag != "" {
		f, err := a.buildFilterForTags(ctx, tag)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode(
				[]repoasproj.RepoAsProj{},
			)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				e.Message {
					Message: "Failed to get repositories",
				},
			)
			prepare("GetProjects", err).Error()
			return
		}
		filter = f
	}

	if name != "" {
		filter = func(f map[string]interface{}) bson.M {
			f["name"] = bson.M{"$regex": name, "$options": "-i"}
			return f
		}(filter)
	}

	if err := a.Repository.Repo.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			c.All(
				context.Background(),
				&repos,
			)
			return c.Err()
		},
		options.Find().
			SetSort(bson.M{"createdat": -1}).
			SetSkip(int64(start)).
			SetLimit(int64(count)),
	); err != mongo.ErrNoDocuments && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to get projects",
			},
		)
		prepare("GetProjects", err).Error()
		return
	}

	projs, err := a.getCompatcProj(repos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to get projects",
			},
		)
		prepare("GetProjects", err).Error()
		return
	}

	sort.Sort(repoasproj.ByCreateDate(projs))

	json.NewEncoder(w).Encode(
		projs,
	)
}

// GetProject
// 
// @Summary return project according to id
// 
// @Description return a project according to id value in path
// 
// @Produce json
// 
// @Router /api/v1/projects/{id} [get]
// 
// @Param id path integer true "a uint value of repository id"
// 
// @Success 200 {object} repoasproj.RepoAsProj
// 
// @Failure 404 {object} e.Message
// 
// @Failure 500 {object} e.Message
func (a *Api) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_id := vars["id"]

	id, _ := strconv.ParseUint(_id, 10, 64)

	var repos []repo.Repo
	ctx := context.Background()

	if err := a.Repository.Repo.GetAllFiltered(
		ctx,
		bson.M{"id": id},
		func(c *mongo.Cursor) error {
			if err := c.All(
				ctx,
				&repos,
			); err != nil {
				return err
			}

			if len(repos) == 0 {
				return mongo.ErrNoDocuments
			}

			return nil
		},
		options.Find(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Don't find project",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to find project",
			},
		)
		prepare("GetProject", err).Error()
		return
	}

	project, err := a.getProjs(repos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to get project",
			},
		)
		prepare("GetProject", err).Error()
		return
	} else if len(project) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Failed to get project",
			},
		)
		prepare("GetProject", fmt.Errorf("Len of project == %v", len(project))).Error()
		return
	}

	json.NewEncoder(w).Encode(
		project[0],	
	)
}

func (a *Api) getProjs(reps []repo.Repo) ([]repoasproj.RepoAsProj, error) {
	var projs []repoasproj.RepoAsProj

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	defer cancel()

	errChan := make(chan error, 1)
	projChan := make(chan []repoasproj.RepoAsProj, 1)

	var count uint
	for i := range reps {
		count++
		go func(r *repo.Repo) {
			proj := repoasproj.RepoAsProj{
				Repo: *r,
			}

			if err := a.Repository.Milestone.GetAllFiltered(
				ctx,
				bson.M{"repoid": r.ID},
				func(c *mongo.Cursor) error {
					c.All(
						context.Background(),
						&proj.Milestones,
					)

					if c.Err() != nil {
						return c.Err()
					}

					return nil
				},
				options.Find(),
			); err != mongo.ErrNoDocuments && err != nil {
				cancel()
				errChan <- err
				return
			}

			if err := a.getAssetsForMilestones(ctx, &proj.Milestones); err != nil {
				cancel()
				errChan <- err
				return
			}
			
			var rl realese.RealeseInRepo
			if err := a.Repository.Realese.GetOne(
				ctx,
				bson.M{"repoid": r.ID},
				func(sr *mongo.SingleResult) error {
					return sr.Decode(&rl)
				},
				options.FindOne(),
			); err != mongo.ErrNoDocuments && err != nil {
				cancel()
				errChan <- err
				return
			} else if err != mongo.ErrNoDocuments {
				proj.LastRealese = &rl.Realese
			}

			if err := a.Repository.Tag.GetAllFiltered(
				ctx,
				bson.M{"repo_id": r.ID},
				func(c *mongo.Cursor) error {
					c.All(
						ctx,
						&proj.Tags,
					)
					return c.Err()
				},
				options.Find(),
			); err != mongo.ErrNoDocuments && err != nil {
				cancel()
				errChan <- err
				return
			}

			projChan <- []repoasproj.RepoAsProj{proj}
		}(&reps[i])
	}

	for i := uint(0); i < count; i++ {
		select {
		case p := <- projChan:
			projs = append(projs, p...)
		case <- ctx.Done():
			err := <- errChan
			return nil, err
		}
	}

	return projs, nil
}

func (a *Api) getCompatcProj(repos []repo.Repo) ([]repoasproj.RepoAsProjCompact, error) {
	var projs []repoasproj.RepoAsProjCompact

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	defer cancel()

	errChan := make(chan error, 1)
	projChan := make(chan []repoasproj.RepoAsProjCompact, 1)

	var count uint
	for i := range repos {
		count++
		go func(r *repo.Repo) {
			proj := repoasproj.RepoAsProjCompact{
				Repo: *r,
			}

			if err := a.Repository.Milestone.GetAllFiltered(
				ctx,
				bson.M{"repoid": r.ID},
				func(c *mongo.Cursor) error {
					var mls []milestone.MilestoneInRepo
					c.All(
						context.Background(),
						&mls,
					)

					if c.Err() != nil {
						return c.Err()
					}
					
					var open float64
					var closed float64
					for _, m := range mls {
						if m.OpenIssues != 0 {
							open += float64(m.OpenIssues)
							closed += float64(m.ClosedIssues)
						}
					}

					if closed == 0 {
						proj.Completed = 1
					} else {
						proj.Completed = open/closed
					}

					return nil
				},
				options.Find(),
			); err != mongo.ErrNoDocuments && err != nil {
				cancel()
				errChan <- err
				return
			}

			if err := a.Repository.Tag.GetAllFiltered(
				ctx,
				bson.M{"repo_id": r.ID},
				func(c *mongo.Cursor) error {
					c.All(
						ctx,
						&proj.Tags,
					)
					return c.Err()
				},
				options.Find(),
			); err != mongo.ErrNoDocuments && err != nil {
				cancel()
				errChan <- err
				return
			}

			projChan <- []repoasproj.RepoAsProjCompact{proj}
		}(&repos[i])
	}

	for i := uint(0); i < count; i++ {
		select {
		case p := <- projChan:
			projs = append(projs, p...)
		case <- ctx.Done():
			err := <- errChan
			return nil, err
		}
	}

	return projs, nil
}

func (a *Api) getAssetsForMilestones(ctx context.Context, ms *[]milestone.Milestone) error {
	for i, m := range *ms {
		var e estimate.Estimate
		if err := a.Repository.Estimate.GetOne(
			ctx,
			bson.M{"milestone_id": m.ID},
			func(sr *mongo.SingleResult) error {
				return sr.Decode(&e)
			},
			options.FindOne(),
		); err != mongo.ErrNoDocuments && err != nil {
			return err
		} else if err != mongo.ErrNoDocuments {
			(*ms)[i].Estimate = &e
		}

		var f functask.FuncTask

		if err := a.Repository.FuncTask.GetOne(
			ctx,
			bson.M{"milestone_id": m.ID},
			func(sr *mongo.SingleResult) error {
				return sr.Decode(&f)
			},
			options.FindOne(),
		); err != mongo.ErrNoDocuments && err != nil {
			return err
		} else if err != mongo.ErrNoDocuments {
			(*ms)[i].FuncTask = &f
		}
	}

	return nil
}

func (a *Api) buildFilterForTags(ctx context.Context, t string) (bson.M, error) {
	var tags []tag.Tag
	massOfTags := strings.Split(t, " ")
	if err := a.Repository.Tag.GetAllFiltered(
		ctx,
		bson.M{"tag": bson.M{"$in": massOfTags}},
		func(c *mongo.Cursor) error {
			if c.RemainingBatchLength() == 0 {
				return mongo.ErrNoDocuments
			}

			return c.All(
				ctx,
				&tags,
			)
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	var ids []uint64
	for _, t := range tags {
		ids = append(ids, t.RepoID)
	}

	return bson.M{"id": bson.M{"$in": ids}}, nil 
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