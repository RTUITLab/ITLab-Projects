package estimate

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddEstimate			endpoint.Endpoint
	DeleteEstimate		endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		AddEstimate: makeAddEstimateEndPoint(s),
		DeleteEstimate: makeDeleteEstimateEndPoint(s),
	}
}

func makeAddEstimateEndPoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*AddEstimateReq)
		err = s.AddEstimate(
			ctx,
			req.EstimateFile,
		)
		if err != nil {
			return nil, err
		}
		return &AddEstimateResp{}, nil
	}
}

func makeDeleteEstimateEndPoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*DeleteEstimateReq)
		err = s.DeleteEstimate(
			ctx,
			req.MilestoneID,
			req.Req,
		)
		if err != nil {
			return nil, err
		}
		return &DeleteEstimateResp{}, nil
	}
}