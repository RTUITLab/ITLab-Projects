package estimate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ITLab-Projects/pkg/statuscode"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ITLab-Projects/service/api/v1/encoder"

	"github.com/gorilla/mux"
)

type AddEstimateReq struct {
	MilestoneID		uint64				`json:"-" swaggerignore:"true"`
	FileID			primitive.ObjectID	`json:"file_id"`
}

type AddEstimateResp struct {
}

func (r *AddEstimateResp) StatusCode() int {
	return http.StatusCreated
}

func (r *AddEstimateResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (r *AddEstimateResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {
	
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

func (r *DeleteEstimateResp) Headers(
	ctx 	context.Context,
	w 		http.ResponseWriter,
) {

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

	vars := mux.Vars(r)

	_mid := vars["milestone_id"]

	milestoneID, _ := strconv.ParseUint(_mid, 10, 64)

	req.MilestoneID = milestoneID

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