package issues

import (
	serverbefore "github.com/ITLab-Projects/service/serverbefore/http"
	"context"

	"github.com/ITLab-Projects/service/api/v1/encoder"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(
	ctx				context.Context,
	endpoints		Endpoints,
	r				*mux.Router,
) {
	r.Handle(
		"/issues",
		httptransport.NewServer(
			endpoints.GetIssues,
			decodeGetIssuesReq,
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