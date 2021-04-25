package v1

import (
	"github.com/ITLab-Projects/pkg/models/functask"
	e "github.com/ITLab-Projects/pkg/err"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"context"
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