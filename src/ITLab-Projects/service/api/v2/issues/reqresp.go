package issues

import (
	"encoding/json"
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/service/api/v1/issues"
)

type GetIssuesReq struct {
	*issues.GetIssuesReq
	HttpURL	*url.URL
}

func decodeGetIssuesReq(
	ctx		context.Context,
	r		*http.Request,
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
		GetIssuesReq: &issues.GetIssuesReq{
			Start: start,
			Count: count,
			Name: name,
			Tag: tag,
		},
		HttpURL: r.URL,
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