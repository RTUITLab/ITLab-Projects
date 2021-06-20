package functask

import (
	"context"

	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddFuncTask			endpoint.Endpoint
	DeleteFuncTask		endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		AddFuncTask: makeAddFuncTaskEndPoint(s),
		DeleteFuncTask: makeDeleteFuncTaskEndPoint(s),
	}
}

func makeAddFuncTaskEndPoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*AddFuncTaskReq)
		err = s.AddFuncTask(
			ctx,
			&functask.FuncTaskFile{
				MilestoneFile: milestonefile.MilestoneFile{
					MilestoneID: req.MilestoneID,
					FileID: req.FileID,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return &AddFuncTaskResp{}, err
	}
}

func makeDeleteFuncTaskEndPoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*DeleteFuncTaskReq)
		err = s.DeleteFuncTask(
			ctx,
			req.MilestoneID,
			req.Req,
		)
		if err != nil {
			return nil, err
		}
		return &DeleteFuncTaskResp{}, nil
	}
}