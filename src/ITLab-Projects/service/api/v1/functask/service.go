package functask

import (
	"context"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/functask"
)

type Service interface {
	AddFuncTask(
		ctx 		context.Context,
		FuncTask	*functask.FuncTaskFile,
	) error

	DeleteFuncTask(
		ctx 		context.Context,
		MilestoneID	uint64,
		// For MFS request
		r			*http.Request,
	) error
}