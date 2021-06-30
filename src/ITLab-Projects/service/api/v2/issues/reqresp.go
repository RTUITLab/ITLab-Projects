package issues

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/urlvalue/encode"
	"github.com/ITLab-Projects/service/api/v1/issues"
)

type GetIssuesReq struct {
	Query GetIssuesQuery
	HttpURL	*url.URL
}

type GetIssuesQuery struct {
	issues.GetIssuesQuery
}

func decodeGetIssuesReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	req := &GetIssuesReq{
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

type GetIssuesResp struct {
	Issues []*milestone.IssuesWithMilestoneID		`json:"items"`
	*chunkresp.ChunkResp							`json:",inline"`
}

func (r *GetIssuesResp) Encode(w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}

func (r *GetIssuesResp) StatusCode() int {
	return http.StatusOK
}

func (r *GetIssuesResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}