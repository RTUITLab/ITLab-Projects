package issues

import (
	"context"

	"github.com/ITLab-Projects/service/api/v1/encoder"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

const (
	GetIssuesName		string	= "get_issues"
	GetLabelsName		string	= "get_labels"
)

func NewHTTPServer(
	ctx			context.Context,
	endpoints	Endpoints,
	r			*mux.Router,
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
		),
	).Methods("GET").Name(GetIssuesName)


	r.Handle(
		"/issues/labels",
		httptransport.NewServer(
			endpoints.GetLabels,
			decodeGetLabels,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("GET").Name(GetLabelsName)
}