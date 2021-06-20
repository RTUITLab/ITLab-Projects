package estimate

import (
	e "github.com/ITLab-Projects/pkg/err"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITLab-Projects/pkg/statuscode"

	"github.com/ITLab-Projects/service/api/v1/beforedelete"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
)

var (
	ErrNotFoundMilestone 		= errors.New("Don't find milestone with this id")
	ErrFailedToSave				= errors.New("Failed to save estimate")
	ErrNotFoundEstimate			= errors.New("Don't find estimate")
	ErrFailedToDeleteEstimate	= errors.New("Failed to delete estimate")
	ErrFileIDIsZero				= errors.New("FileID can't be zero")
)

func init() {
	// to generate swagger
	_ = e.Message{}
}

type service struct {
	repository 	Repository
	logger		log.Logger
	mfsreq		mfsreq.Requester
}

func New(
	Repository 	Repository,
	logger 		log.Logger,
	MFSReq		mfsreq.Requester,
) *service {
	return &service{
		repository: Repository,
		logger: logger,
		mfsreq: MFSReq,
	}
}

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
// @Router /v1/estimate/{milestone_id} [post]
//
// @Security ApiKeyAuth
// 
// @Accept json
//
// @Produce json
//
// @Param estimate body estimate.AddEstimateReq true "estimate that you want to add"
// 
// @Param milestone_id path integer true "id of milestone"
//
// @Success 201
//
// @Failure 400 {object} e.Message "Unexpected body"
//
// @Failure 500 {object} e.Message "Failed to save estimate"
//
// @Failure 404 {object} e.Message "Don't find milestone with this id"
// 
// @Failure 401 {object} e.Message 
// 
// @Failure 403 {object} e.Message "if you are not admin"
func (s *service) AddEstimate(
	ctx context.Context, 
	est *estimate.EstimateFile,
) error {
	logger := log.With(s.logger,"method", "AddEstimate")
	if est.FileID.IsZero() {
		return statuscode.WrapStatusError(
			ErrFileIDIsZero,
			http.StatusBadRequest,
		)
	}
	// Check if milestone with this id exists
	_, err := s.repository.GetMilestoneByID(
		ctx,
		est.MilestoneID,
	); 
	switch {
	case err == mongo.ErrNoDocuments:
		return statuscode.WrapStatusError(
			ErrNotFoundMilestone,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to save estimate: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToSave,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveEstimate(
		ctx,
		est,
	); err != nil {
		level.Error(logger).Log("Failed to save estimate: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToSave,
			http.StatusInternalServerError,
		)
	}

	return nil
}

// DeleteEstimate
// 
// @Tags estimate
// 
// @Summary delete estimate from database
// 
// @Description delete estimate from database
// 
// @Security ApiKeyAuth
// 
// @Router /v1/estimate/{milestone_id} [delete]
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
// 
// @Failure 409 {object} e.Message "some problems with microfileservice"
// 
// @Failure 401 {object} e.Message 
// 
// @Failure 403 {object} e.Message "if you are not admin"
func (s *service) DeleteEstimate(
	ctx context.Context,
	MilestoneID uint64, 
	r *http.Request,
) error {
	logger := log.With(s.logger,"method", "DeleteEstimate")
	est, err := s.repository.GetEstimateByMilestoneID(
		ctx,
		MilestoneID,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return statuscode.WrapStatusError(
			ErrNotFoundEstimate,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete estimate: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteEstimate,
			http.StatusInternalServerError,
		)
	}

	err = beforedelete.BeforeDeleteWithReq(
		s.mfsreq,
		r,
	)(est)
	switch {
	case errors.Is(err, mfsreq.NetError):
		level.Error(logger).Log("Failed to delete estimate: err", err)
		return statuscode.WrapStatusError(
			mfsreq.NetError,
			http.StatusConflict,
		)
	case mfsreq.IfUnexcpectedCode(err):
		uce := err.(*mfsreq.UnexpectedCodeErr)
		causedErr := fmt.Errorf("Unecxpected code from microfileserver: %v", uce.Code)
		level.Error(logger).Log("Failed to delete estimate: err", causedErr)
		return statuscode.WrapStatusError(
			causedErr,
			http.StatusConflict,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete estimate: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteEstimate,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.DeleteOneEstimateByMilestoneID(
		ctx,
		MilestoneID,
	); err != nil {
		level.Error(logger).Log("Failed to delete estimate: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteEstimate,
			http.StatusInternalServerError,
		)
	}

	return nil
}

