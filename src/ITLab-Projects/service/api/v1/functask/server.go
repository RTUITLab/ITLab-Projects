package functask

import (
	"context"

	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

const (
	AddFuncTaskName		string 	= "add_functask_admin"
	DeleteFuncTaskName	string	= "delete_functasl_admin"
)

func NewHTTPServer(
	ctx			context.Context,
	endpoints	Endpoints,
	r			*mux.Router,
) {
	r.Handle(
		"/task",
		httptransport.NewServer(
			endpoints.AddFuncTask,
			decodeAddFuncTaskReq,
			encodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("POST").Name(AddFuncTaskName)

	r.Handle(
		"/task/{milestone_id:[0-9]+}",
		httptransport.NewServer(
			endpoints.DeleteFuncTask,
			decodeDeleteFuncTaskReq,
			encodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("DELETE").Name(DeleteFuncTaskName)
}