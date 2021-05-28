package tags

import (
	"context"

	"github.com/ITLab-Projects/service/api/v1/encoder"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

const (
	GetAllTagsName		string = "get_all_tags"
)

func NewHTTPServer(
	ctx 		context.Context,
	endpoints 	Endpoints,
	r			*mux.Router,
) {
	r.Handle(
		"/tags",
		httptransport.NewServer(
			endpoints.GetAllTags,
			decodeGetAllTagsReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("GET").Name(GetAllTagsName)
}