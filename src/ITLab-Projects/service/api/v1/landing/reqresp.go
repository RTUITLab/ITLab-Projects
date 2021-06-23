package landing

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/urlvalue/encode"
	"github.com/gorilla/mux"
)

type GetAllLandingsReq struct {
	Query	GetAllLandingsQuery
}

type GetAllLandingsQuery struct {
	Start 	int64	`query:"start,int"`
	Count 	int64	`query:"count,int"`
	Tag		string	`query:"tag,string"`
	Name	string	`query:"name,string"`
}

func decodeGetAllLandingsReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	req := &GetAllLandingsReq{}

	if err := encode.UrlQueryUnmarshall(
		&req.Query,
		r.URL.Query(),
	); err != nil {
		return nil, err
	}

	return req, nil
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