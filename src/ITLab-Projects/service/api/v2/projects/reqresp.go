package projects

import (
	"strconv"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/ITLab-Projects/service/api/v1/projects"
)

type GetProjectsReq struct {
	*projects.GetProjectsReq
	HttpURL		*url.URL
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
		GetProjectsReq: &projects.GetProjectsReq{
			Start: start,
			Count: count,
			Name: name,
			Tag: tag,
		},
		HttpURL: r.URL,
	}, nil
}

type GetProjectsResp struct {
	Projects []*repoasproj.RepoAsProjCompactPointers		`json:"items"`
	*chunkresp.ChunkResp									`json:",inline"`
}

func (r *GetProjectsResp) Encode(w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
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