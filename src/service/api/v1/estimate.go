package v1

import (
	"github.com/ITLab-Projects/pkg/models/estimate"
	e "github.com/ITLab-Projects/pkg/err"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"context"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

// AddEstimate
// 
// @Tags estimate
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
// @Tags estimate
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