package v1

import (
	"github.com/ITLab-Projects/pkg/mfsreq"
	"sort"
	"strconv"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"strings"
	"github.com/ITLab-Projects/pkg/models/functask"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"errors"
	"github.com/ITLab-Projects/pkg/models/repo"
	_ "time"
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"encoding/json"
	"github.com/ITLab-Projects/pkg/githubreq"
	e "github.com/ITLab-Projects/pkg/err"
	"context"
	"net/http"
)

// UpdateAllProjects
// 
// @Tags projects
// 
// @Summary Update all projects
// 
// @Description make all request to github to update repositories, milestones
// 
// @Router /api/v1/projects/ [post]
// 
// @Success 200
// 
// @Failure 409 {object} e.Err
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
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(
			e.Err {
				Err: err.Error(),
				Message: e.Message{
					Message: "Failed to update",
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
		// TODO delete in prod
		// time.Sleep(1*time.Second)
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
		// TODO delete in prod
		// time.Sleep(2*time.Second)
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
		// TODO delete in prod
		// time.Sleep(3*time.Second)
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
			w.WriteHeader(http.StatusConflict)
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

	if err := a.Repository.Repo.SaveAndUpdatenUnfind(
		context.Background(),
		repo.ToRepo(repos),
		bson.M{"$set": bson.M{"deleted": true}},
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e.Message{
			Message: "Can't save repositories",
		})
		prepare("UpdateAllProjects", err).Error("Can't save repositories")
		return
	}

	if err := a.Repository.Milestone.SaveAndUpdatenUnfind(
		context.Background(),
		ms,
		bson.M{"$set": bson.M{"deleted": true}},
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
}

// GetProjects
// 
// @Tags projects
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
		// "createdat": -1, "deleted": -1,
			SetSort(bson.D{ {"createdat", -1}, {"deleted", 1}} ).
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
// @Tags projects
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

// DeleteProjects
// 
// @Summary delete project by id
// 
// @Description delete project by id and all milestones in it
// 
// @Tags projects
// 
// @Router /api/v1/projects/{id} [delete]
// 
// @Param id path integer true "id of project"
// 
// @Success 200
// 
// @Failure 404 {object} e.Message
// 
// @Failure 500 {object} e.Message
// 
// @Failure 409 {object} e.Message "some problems with microfileservice"
func (a *Api) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_id := vars["id"]

	id, _ := strconv.ParseUint(_id, 10, 64)

	var rep repo.Repo

	if err := a.Repository.Repo.GetOne(
		context.Background(),
		bson.M{"id": id},
		func(sr *mongo.SingleResult) error {
			if err :=  sr.Decode(&rep); err != nil {
				return err
			}

			return sr.Err()
		},
		options.FindOne(),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Project not found",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Can't delete project",
			},
		)
		prepare("DeleteProject", err).Error()
		return
	}

	if err := a.deleteMilestones(
		context.Background(),
		rep.ID,
		a.beforeDelete,
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if errors.Is(err, mfsreq.NetError) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Faield to delete project",
			},
		)
		prepare("DeleteProject", err).Error()
		return
	} else if errors.Is(err, mfsreq.ErrUnexpectedCode) {
		uce := err.(*mfsreq.UnexpectedCodeErr)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(
				e.Message {
					Message: fmt.Sprintf("Unecxpected code: %v", uce.Code),
				},
			)
			prepare("DeleteProject", err).Error()
			return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Failed to delete project",
			},
		)
		prepare("DeleteProject", err).Error()
		return
	}

	if err := a.Repository.Repo.DeleteOne(
		context.Background(),
		bson.M{"id": id},
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
			e.Message {
				Message: "Project not found",
			},
		)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "Can't delete project",
			},
		)
		prepare("DeleteProject", err).Error()
		return
	}


}

// Return deleted milestone and error
func (a *Api) deleteMilestones(
	ctx context.Context, 
	repoid uint64,
	// switch type of argument to estimate and functtask
	beforeDelete beforeDeleteFunc,
) (error) {
	var ms []milestone.Milestone

	if err := a.Repository.Milestone.GetAllFiltered(
		ctx,
		bson.M{"repoid": repoid},
		func(c *mongo.Cursor) error {
			if err := c.All(
				ctx,
				&ms,
			); err != nil {
				return err
			}

			if err := c.Err(); err != nil {
				return err
			}

			if ms == nil {
				return mongo.ErrNoDocuments
			}

			return nil
		},
		options.Find(),
	); err != nil {
		return err
	}

	if err := beforeDelete(ms); err != nil {
		return err
	}

	if err := a.deleteEstimates(
		ctx,
		func(ms []milestone.Milestone) []uint64 {
			var ids []uint64

			for _, m := range ms {
				ids = append(ids, m.ID)
			}

			return ids
		}(ms),
		beforeDelete,
	); err != nil {
		return err
	}

	if err := a.deleteFuncTasks(
		ctx,
		func(ms []milestone.Milestone) []uint64 {
			var ids []uint64

			for _, m := range ms {
				ids = append(ids, m.ID)
			}

			return ids
		}(ms),
		beforeDelete,
	); err != nil {
		return err
	}

	if err := a.Repository.Milestone.DeleteMany(
		ctx,
		bson.M{"repoid": repoid},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}

			return nil
		},
		options.Delete(),
	); err != nil {
		return err
	}

	return nil
}

func (a *Api) deleteTags(ctx context.Context, repoid uint64) error {
	if err := a.Repository.Tag.DeleteMany(
		ctx,
		bson.M{"repoi_d": repoid},
		nil,
		options.Delete(),
	); err != nil {
		return err
	}

	return nil
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
		var e estimate.EstimateFile
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
			(*ms)[i].Estimate = &estimate.Estimate{
				MilestoneID: e.MilestoneID,
				EstimateURL: a.MFSRequester.GenerateDownloadLink(e.FileID),
			}
		}

		var f functask.FuncTaskFile

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
			(*ms)[i].FuncTask = &functask.FuncTask{
				MilestoneID: f.MilestoneID,
				FuncTaskURL: a.MFSRequester.GenerateDownloadLink(f.FileID),
			}
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