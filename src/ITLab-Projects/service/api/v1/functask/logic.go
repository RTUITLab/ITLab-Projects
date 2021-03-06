package functask

import (
	e "github.com/ITLab-Projects/pkg/err"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"context"
	"errors"
	"net/http"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/ITLab-Projects/service/api/v1/beforedelete"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	_ = e.Err{}
}

var (
	ErrNotFoundMilestone 		= errors.New("Don't find milestone with this id")
	ErrFailedToSave				= errors.New("Failed to save functask")
	ErrNotFound					= errors.New("Don't find functask")
	ErrFailedToDelete			= errors.New("Failed to delete functask")
	ErrFileIDNil				= errors.New("File id can't be nil")
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
// @Router /v1/task [post]
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
// @Failure 403 {object} e.Message "if you are not admin"
func (s *service) AddFuncTask(
	ctx 		context.Context, 
	FuncTask 	*functask.FuncTaskFile,
) error {
	logger := log.With(s.logger, "method", "AddFuncTask")
	if FuncTask.FileID.IsZero() {
		return statuscode.WrapStatusError(
			ErrFileIDNil,
			http.StatusBadRequest,
		)
	}
	// Check if milestone with this id exists
	_, err := s.repository.GetMilestoneByID(
		ctx,
		FuncTask.MilestoneID,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return statuscode.WrapStatusError(
			ErrNotFoundMilestone,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to save functask: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToSave,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.SaveFuncTask(
		ctx,
		FuncTask,
	); err != nil {
		level.Error(logger).Log("Failed to save functask: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToSave,
			http.StatusInternalServerError,
		)
	}

	return nil
}

// DeleteFuncTask
// 
// @Tags functask
// 
// @Summary delete functask from database
// 
// @Description delete functask from database
// 
// @Router /v1/task/{milestone_id} [delete]
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
// @Failure 403 {object} e.Message "if you are not admin"
func (s *service) DeleteFuncTask(
	ctx 		context.Context, 
	MilestoneID uint64, 
	r 			*http.Request,
) error {
	logger := log.With(s.logger, "method", "DeleteFuncTask")
	ft, err := s.repository.GetFuncTaskByMilestoneID(
		ctx,
		MilestoneID,
	)
	switch {
	case err == mongo.ErrNoDocuments:
		return statuscode.WrapStatusError(
			ErrNotFound,
			http.StatusNotFound,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete functask: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDelete,
			http.StatusInternalServerError,
		)
	}

	err = beforedelete.BeforeDeleteWithReq(
		s.mfsReq,
		r,
	)(ft)
	switch {
	case errors.Is(err, mfsreq.NetError):
		level.Error(logger).Log("Failed to delete functask: err", err)
		return statuscode.WrapStatusError(
			mfsreq.NetError,
			http.StatusConflict,
		)
	case mfsreq.IfUnexcpectedCode(err):
		uce := err.(*mfsreq.UnexpectedCodeErr)
		causedErr := fmt.Errorf("Unecxpected code from microfileserver: %v", uce.Code)
		level.Error(logger).Log("Failed to delete funcTask: err", causedErr)
		return statuscode.WrapStatusError(
			causedErr,
			http.StatusConflict,
		)
	case err != nil:
		level.Error(logger).Log("Failed to delete functask: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDelete,
			http.StatusInternalServerError,
		)
	}

	if err := s.repository.DeleteOneFuncTaskByMilestoneID(
		ctx,
		MilestoneID,
	); err != nil {
		level.Error(logger).Log("Failed to delete functask: err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDelete,
			http.StatusInternalServerError,
		)
	}

	return nil
}

