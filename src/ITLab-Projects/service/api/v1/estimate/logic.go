package estimate

import (
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
)

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

func (s *service) AddEstimate(
	ctx context.Context, 
	est *estimate.EstimateFile,
) error {
	logger := log.With(s.logger,"method", "AddEstimate")
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

