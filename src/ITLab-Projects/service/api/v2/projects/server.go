package projects

import (
	serverbefore "github.com/ITLab-Projects/service/serverbefore/http"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	"github.com/ITLab-Projects/service/api/v1/encoder"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"context"
)

func NewHTTPServer(
	ctx			context.Context,
	endpoints	Endpoints,
	r			*mux.Router,
) {
	r.Handle(
		"/projects",
		httptransport.NewServer(
			endpoints.GetProjects,
			decodeGetProjectsReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
		),
	).Methods("GET")
}