package tags

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/tag"
)

type GetAllTagsReq struct {
}

type GetAllTagsResp struct {
	Tags []*tag.Tag
}

func (r *GetAllTagsResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Tags)
}

func (r *GetAllTagsResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}

func (r *GetAllTagsResp) StatusCode() int {
	return http.StatusOK
}

func decodeGetAllTagsReq(
	ctx 	context.Context,
	r		*http.Request,
) (interface{}, error) {
	return &GetAllTagsReq{}, nil
}