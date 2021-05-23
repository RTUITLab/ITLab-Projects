package estimate

import (
	"fmt"
	"github.com/ITLab-Projects/pkg/statuscode"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/service/api/v1/encoder"

	"github.com/gorilla/mux"

	"github.com/ITLab-Projects/pkg/models/estimate"
)

type AddEstimateReq struct {
	*estimate.EstimateFile
}

type AddEstimateResp struct {
}

func (r *AddEstimateResp) StatusCode() int {
	return http.StatusCreated
}

func (r *AddEstimateResp) Encode(w http.ResponseWriter) error {
	return nil
}

type DeleteEstimateReq struct {
	MilestoneID		uint64
	Req				*http.Request
}

type DeleteEstimateResp struct {
}

func (r *DeleteEstimateResp) StatusCode() int {
	return http.StatusOK
}

func (r *DeleteEstimateResp) Encode(w http.ResponseWriter) error {
	return nil
}

func decodeAddEstimateReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	req := &AddEstimateReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, statuscode.WrapStatusError(
			fmt.Errorf("Unexcpected body"),
			http.StatusBadRequest,
		)
	}
	return req, nil
}

func decodeDeleteEstimateReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	req := &DeleteEstimateReq{}
	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	req.MilestoneID = milestoneID
	req.Req = r

	return req, nil
}

func encodeResponce(
	ctx context.Context, 
	w http.ResponseWriter, 
	resp interface{},
) error {
	return encoder.EncodeResponce(ctx,w,resp)
}