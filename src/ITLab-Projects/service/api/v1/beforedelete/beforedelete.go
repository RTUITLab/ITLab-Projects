package beforedelete

import (
	"net/http"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
)

type BeforeDeleteFunc func(interface{}) error

func BeforeDeleteWithReq(
	MFSReq mfsreq.Requester,
	r *http.Request,
) BeforeDeleteFunc {
	return func(v interface{}) error {
		return BeforeDelete(
			MFSReq.NewRequests(r),
			v,
		)
	}
}

func BeforeDelete(
	deleter mfsreq.FileDeleter,
	v interface{},
) error {
	switch v.(type) {
	case estimate.EstimateFile:
		est, _ := v.(estimate.EstimateFile)
		if err := deleter.DeleteFile(est.FileID); err != nil {
			return err
		}
	case []estimate.EstimateFile:
		ests, _ := v.([]estimate.EstimateFile)
		for _, est := range ests {
			if err := deleter.DeleteFile(est.FileID); err != nil {
				return err
			}
		}
	case []*estimate.EstimateFile:
		ests, _ := v.([]*estimate.EstimateFile)
		for _, est := range ests {
			if err := deleter.DeleteFile(est.FileID); err != nil {
				return err
			}
		}
	case functask.FuncTaskFile:
		task, _ := v.(functask.FuncTaskFile)
		if err := deleter.DeleteFile(task.FileID); err != nil {
			return err
		}
	case []functask.FuncTaskFile:
		tasks, _ := v.([]functask.FuncTaskFile)
		for _, task := range tasks {
			if err := deleter.DeleteFile(task.FileID); err != nil {
				return err
			}
		}
	case []*functask.FuncTaskFile:
		tasks, _ := v.([]*functask.FuncTaskFile)
		for _, task := range tasks {
			if err := deleter.DeleteFile(task.FileID); err != nil {
				return err
			}
		}
	default:
	}

	return nil
}