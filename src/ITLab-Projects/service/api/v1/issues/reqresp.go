package issues

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/urlvalue/encode"
)

type GetIssuesQuery struct {
	Start	int		`query:"start,int"`
	Count	int		`query:"count,int"`
	Name	string	`query:"name,string"`
	Tag		string	`query:"tag,string"`
}

type GetIssuesReq struct {
	Query	GetIssuesQuery
}

type GetIssuesResp struct {
	Issues []*milestone.IssuesWithMilestoneID
}

func (r *GetIssuesResp) StatusCode() int {
	return http.StatusOK
}

func (r *GetIssuesResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Issues)
}

func (r *GetIssuesResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

type GetLabelsReq struct {

}

type GetLabelsResp struct {
	Labels	[]interface{}
}

func (r *GetLabelsResp) StatusCode() int {
	return http.StatusOK
}

func (r *GetLabelsResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Labels)
}

func (r *GetLabelsResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func decodeGetIssuesReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	values := r.URL.Query()

	req := &GetIssuesReq{
	}

	encode.UrlQueryUnmarshall(
		&req.Query,
		values,
	)

	return req, nil
}

func decodeGetLabels(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	return &GetLabelsReq{}, nil
}