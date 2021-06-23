package projects

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/ITLab-Projects/pkg/urlvalue/encode"
	"github.com/ITLab-Projects/service/api/v1/projects"
)

type GetProjectsReq struct {
	Query		GetProjectsQuery
	HttpURL		*url.URL
}

type GetProjectsQuery struct {
	projects.GetProjectsQuery
}

func decodeGetProjectsReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	req := &GetProjectsReq{
		HttpURL: r.URL,
	}

	if err := encode.UrlQueryUnmarshall(
		&req.Query,
		r.URL.Query(),
	); err != nil {
		return nil, err
	}

	return req, nil
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