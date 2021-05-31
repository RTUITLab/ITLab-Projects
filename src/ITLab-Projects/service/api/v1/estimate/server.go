package estimate

import (
	serverbefore "github.com/ITLab-Projects/service/serverbefore/http"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"
	"context"

	"github.com/gorilla/mux"
)

const (
	AddEstimateName		string 	= "add_estimate_admin"
	DeleteEstimateName 	string	= "delete_estimate_admin" 
)

// Make http endpoint
// 
// To add middleware use mux.WalkFunc
func NewHTTPServer(
	ctx 		context.Context,
	endpoints	Endpoints,
	r			*mux.Router,
) {
	r.Handle(
		"/estimate",
		httptransport.NewServer(
			endpoints.AddEstimate,
			decodeAddEstimateReq,
			encodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
		),
	).Methods("POST").Name(AddEstimateName)

	r.Handle(
		"/estimate/{milestone_id:[0-9]+}",
		httptransport.NewServer(
			endpoints.DeleteEstimate,
			decodeDeleteEstimateReq,
			encodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
		),
	).Methods("DELETE").Name(DeleteEstimateName)
}