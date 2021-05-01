package v1

import (
	"fmt"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"errors"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddFuncTask
//
// @Tags functask
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
// @Param functask body functask.FuncTaskFile true "function task that you want to add"
//
// @Success 201
//
// @Failure 400 {object} e.Message "Unexpected body"
//
// @Failure 500 {object} e.Message "Failed to save functask"
//
// @Failure 404 {object} e.Message "Don't find milestone with this id"
// @Failure 401 {object} e.Message 
// 
// @Failure 403 {object} e.Message "if you are nor admin"
func (a *Api) AddFuncTask(w http.ResponseWriter, r *http.Request) {
	var fntask functask.FuncTaskFile
	if err := json.NewDecoder(r.Body).Decode(&fntask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Unexpected body",
			},
		)
		return
	}

	if fntask.FileID.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			e.Message {
				Message: "id is not an objectid",
			},
		)
		prepare("AddFuncTask", fmt.Errorf("id is not an objectid")).Warn()
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
// @Tags functask
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
// 
// @Failure 409 {object} e.Message "some problems with microfileservice"
// 
// @Failure 401 {object} e.Message 
// 
// @Failure 403 {object} e.Message "if you are nor admin"
func (a *Api) DeleteFuncTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10 * time.Second,
	)
	defer cancel()

	if err := a.deleteFuncTask(
		ctx,
		milestoneID,
		a.beforeDeleteWithReq(r),
	); err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "func task not found",
			},
		)
		return
	} else if errors.Is(err, mfsreq.NetError) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(
			e.Message{
				Message: "Faield to delete functask",
			},
		)
		prepare("DeleteFuncTask", err).Error()
		return 
	} else if errors.Is(err, mfsreq.ErrUnexpectedCode) {
		uce := err.(*mfsreq.UnexpectedCodeErr)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(
				e.Message {
					Message: fmt.Sprintf("Unecxpected code: %v", uce.Code),
				},
			)
			prepare("DeleteFuncTask", err).Error()
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

func (a *Api) deleteFuncTask(
	ctx context.Context, 
	milestoneid uint64,
	beforeDelete beforeDeleteFunc,
	) (error) {
	var task functask.FuncTaskFile

	if err := a.Repository.FuncTask.GetOne(
		ctx,
		bson.M{"milestone_id": milestoneid},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&task)
		},
		options.FindOne(),
	); err != nil {
		return err
	}

	if err := beforeDelete(task); err != nil {
		return err
	}

	if err := a.Repository.FuncTask.DeleteOne(
		ctx,
		bson.M{"milestone_id": milestoneid},
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

// Ignore 404 status code
func (a *Api) deleteFuncTasks(
	ctx context.Context, 
	milestonesid []uint64,
	beforeDelete beforeDeleteFunc,
	) (error) {
	var tasks []functask.FuncTaskFile

	if err := a.Repository.FuncTask.GetAllFiltered(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": milestonesid}},
		func(c *mongo.Cursor) error {
			if err := c.All(
				ctx,
				&tasks,
			); err != nil {
				return err
			}

			return c.Err()
		},
		options.Find(),
	); err != nil {
		return err
	}

	if err := beforeDelete(tasks); err != nil {
		return err
	}

	if err := a.Repository.FuncTask.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": milestonesid}},
		nil,
		options.Delete(),
	); err != nil {
		return err
	}

	return nil
}