package tags

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAllTags		endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetAllTags: makeGetAllTagsEndpoint(s),
	}
}

func makeGetAllTagsEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		_ = request.(*GetAllTagsReq)
		tgs, err := s.GetAllTags(
			ctx,
		)
		if err != nil {
			return nil, err
		}

		return &GetAllTagsResp{Tags: tgs}, nil
	}
}