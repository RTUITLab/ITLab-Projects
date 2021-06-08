package landing

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/gorilla/mux"
)

type GetAllLandingsReq struct {
	Start 	int64
	Count 	int64
	Tag		string
	Name	string
}

func decodeGetAllLandingsReq(
	ctx		context.Context,
	r		*http.Request,
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

	return &GetAllLandingsReq{
		Start: start,
		Count: count,
		Tag: tag,
		Name: name,
	}, nil
}

type GetAllLandingResp struct {
	Landings []*landing.LandingCompact
}

func (r *GetAllLandingResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Landings)
}

func (r *GetAllLandingResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}

func (r *GetAllLandingResp) StatusCode() int {
	return http.StatusOK
}

type GetLandingReq struct {
	ID		uint64
}

func decodeGetLandingReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)

	_id := vars["id"]

	ID, _ := strconv.ParseInt(_id, 10, 64)

	return &GetLandingReq{
		ID: uint64(ID),
	}, nil
}

type GetLandingResp struct {
	Landing *landing.Landing
}

func (r *GetLandingResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.Landing)
}

func (r *GetLandingResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	w.Header().Add("Content-Type", "application/json")
}

func (r *GetLandingResp) StatusCode() int {
	return http.StatusOK
}