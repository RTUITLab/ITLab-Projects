package updater

import (
	"net/http"
	"context"
)

type UpdateProjectsReq struct {

}

func decodeUpdateProjetcsReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	return &UpdateProjectsReq{}, nil
}

type UpdateProjectsResp struct {

}

func (r *UpdateProjectsResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (r *UpdateProjectsResp) StatusCode() int {
	return http.StatusOK
}

func (r *UpdateProjectsResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {

}