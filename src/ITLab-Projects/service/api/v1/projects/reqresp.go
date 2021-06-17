package projects

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/gorilla/mux"
)

type GetProjectReq struct {
	ID		uint64
}

func decodeGetProjectReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)
	_id := vars["id"]

	ID, _ := strconv.ParseUint(_id, 10, 64)

	return &GetProjectReq{
		ID: ID,
	}, nil
}

type GetProjectResp struct {
	*repoasproj.RepoAsProjPointer
}

func (r *GetProjectResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.RepoAsProjPointer)
}

func (r *GetProjectResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}

func (r *GetProjectResp) StatusCode() int {
	return http.StatusOK
}

type GetProjectsReq struct {
	Start, 	Count 	int64
	Name, 	Tag		string
}

func decodeGetProjectsReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	values := r.URL.Query()

	_start := values.Get("start")
	_count := values.Get("count")
	name := values.Get("name")
	tag := values.Get("tag")

	start, err := strconv.ParseInt(_start, 10, 64)
	if err != nil {
		start = 0
	}

	count, err := strconv.ParseInt(_count, 10, 64)
	if err != nil {
		count = 0
	}

	return &GetProjectsReq{
		Start: start,
		Count: count,
		Name: name,
		Tag: tag,
	}, nil
}

type GetProjectsResp struct {
	Projects []*repoasproj.RepoAsProjCompactPointers
}

func (r *GetProjectsResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Projects)
}

func (r *GetProjectsResp) StatusCode() int {
	return http.StatusOK
}

func (r *GetProjectsResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}

type DeleteProjectReq struct {
	ID		uint64
	Req		*http.Request
}

func decodeDeleteProjectsReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)

	_id := vars["id"]
	ID, _ := strconv.ParseUint(_id, 10, 64)

	return &DeleteProjectReq{
		ID: ID,
		Req: r,
	}, nil
}

type DeleteProjectsResp struct {

}

func (r *DeleteProjectsResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (r *DeleteProjectsResp) StatusCode() int {
	return http.StatusOK
}

func (r *DeleteProjectsResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {

}