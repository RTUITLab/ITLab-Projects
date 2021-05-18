package functask

import (
	"context"
	"net/http"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/service/api/v1/beforedelete"
	"github.com/go-kit/kit/log"
)

type service struct {
	repository 	Repository
	mfsReq		mfsreq.Requester
	logger		log.Logger
}

func New(
	Repository 	Repository,
	MFSReq		mfsreq.Requester,
	Logger		log.Logger,
) *service {
	return &service{
		repository: Repository,
		mfsReq: MFSReq,
		logger: Logger,
	}
}

func (s *service) AddFuncTask(
	ctx 		context.Context, 
	FuncTask 	*functask.FuncTaskFile,
) error {
	// Check if milestone with this id exists
	if _, err := s.repository.GetMilestoneByID(
		ctx,
		FuncTask.MilestoneID,
	); err != nil { // should be nil
		return err
	}

	if err := s.repository.SaveFuncTask(
		ctx,
		FuncTask,
	); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteFuncTask(
	ctx 		context.Context, 
	MilestoneID uint64, 
	r 			*http.Request,
) error {
	ft, err := s.repository.GetFuncTaskByMilestoneID(
		ctx,
		MilestoneID,
	)
	if err != nil {
		return err
	}

	if err := beforedelete.BeforeDeleteWithReq(
		s.mfsReq,
		r,
	)(ft); err != nil {
		return err
	}

	if err := s.repository.DeleteOneFuncTaskByMilestoneID(
		ctx,
		MilestoneID,
	); err != nil {
		return err
	}

	return nil
}

