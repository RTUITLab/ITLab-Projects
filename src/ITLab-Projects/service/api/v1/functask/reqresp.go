package functask

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/ITLab-Projects/service/api/v1/encoder"
)

type AddFuncTaskReq struct {
	FileID 			primitive.ObjectID	`json:"file_id"`
	MilestoneID		uint64				`json:"-" swaggerignore:"true"`
}

type AddFuncTaskResp struct {
}

func (r *AddFuncTaskResp) StatusCode() int {
	return	http.StatusCreated
}

func (r *AddFuncTaskResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (r *AddFuncTaskResp) Headers(ctx context.Context, w http.ResponseWriter) {

}

type DeleteFuncTaskReq struct {
	MilestoneID		uint64
	Req				*http.Request
}

type DeleteFuncTaskResp struct {
}

func (r *DeleteFuncTaskResp) StatusCode() int {
	return	http.StatusOK
}

func (r *DeleteFuncTaskResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (r *DeleteFuncTaskResp) Headers(ctx context.Context, w http.ResponseWriter) {

}

func decodeAddFuncTaskReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	req := &AddFuncTaskReq{}

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

func decodeDeleteFuncTaskReq(
	ctx context.Context,
	r	*http.Request,
) (interface{}, error) {
	req := &DeleteFuncTaskReq{}

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
	return encoder.EncodeResponce(ctx, w, resp)
}