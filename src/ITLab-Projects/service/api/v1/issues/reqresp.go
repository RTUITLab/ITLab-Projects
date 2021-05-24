package issues

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/pkg/models/milestone"
)

type GetIssuesReq struct {
	Start, 	Count 		int64
	Name, 	Tag 		string
}

type GetIssuesResp struct {
	Issues []*milestone.IssuesWithMilestoneID
}

func (r *GetIssuesResp) StatusCode() int {
	return http.StatusOK
}

func (r *GetIssuesResp) Encode(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(r.Issues)
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
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(r.Labels)
}

func decodeGetIssuesReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	values := r.URL.Query()

	_start 	:= values.Get("start")
	_count 	:= values.Get("count")
	name 	:= values.Get("name")
	tag		:= values.Get("tag")

	start, err := strconv.ParseInt(_start, 10, 64)
	if err != nil {
		start = 0
	}

	count, err := strconv.ParseInt(_count, 10, 64)
	if err != nil {
		count = 0
	}

	req := &GetIssuesReq{
		Start: start,
		Count: count,
		Name: name,
		Tag: tag,
	}

	return req, nil
}

func decodeGetLabels(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	return &GetLabelsReq{}, nil
}