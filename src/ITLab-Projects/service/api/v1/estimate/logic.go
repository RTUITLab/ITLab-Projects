package estimate

import (
	"context"
	"net/http"

	"github.com/ITLab-Projects/service/api/v1/beforedelete"
	"github.com/go-kit/kit/log"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
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
	// Check if milestone with this id exists
	if _, err := s.repository.GetMilestoneByID(
		ctx,
		est.MilestoneID,
	); err != nil { // should be nil
		return err
	}

	if err := s.repository.SaveEstimate(
		ctx,
		est,
	); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteEstimate(
	ctx context.Context,
	MilestoneID uint64, 
	r *http.Request,
) error {
	est, err := s.repository.GetEstimateByMilestoneID(
		ctx,
		MilestoneID,
	)
	if err != nil { // Should be nil
		return err
	}

	if err := beforedelete.BeforeDeleteWithReq(
		s.mfsreq,
		r,
	)(est); err != nil {
		return err
	}

	if err := s.repository.DeleteOneEstimateByMilestoneID(
		ctx,
		MilestoneID,
	); err != nil {
		return err
	}

	return nil
}

