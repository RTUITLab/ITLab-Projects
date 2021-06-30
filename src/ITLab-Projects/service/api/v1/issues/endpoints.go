package issues

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetIssues	endpoint.Endpoint
	GetLabels	endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		GetIssues: makeGetIssuesEndpoint(s),
		GetLabels: makeGetlabelsEndpoint(s),
	}
}

func makeGetIssuesEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetIssuesReq)
		is, err := s.GetIssues(
			ctx,
			req.Query,
		)
		if err != nil {
			return nil, err
		}

		return &GetIssuesResp{
			Issues: is,
		}, nil
	}
}

func makeGetlabelsEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		_ = request.(*GetLabelsReq)

		labels, err := s.GetLabels(
			ctx,
		)
		if err != nil {
			return nil, err
		}

		return &GetLabelsResp{
			Labels: labels,
		}, nil
	}
}