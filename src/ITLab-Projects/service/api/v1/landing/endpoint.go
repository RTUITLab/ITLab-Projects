package landing

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetLanding		endpoint.Endpoint
	GetAllLandings	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetLanding: makeGetLandingEndpoint(s),
		GetAllLandings: makeGetLandingsEndpoint(s),
	}
}

func makeGetLandingEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetLandingReq)
		landing, err := s.GetLanding(
			ctx,
			req.ID,
		)
		if err != nil {
			return nil, err
		}

		return &GetLandingResp{
			Landing: landing,
		}, nil
	}
}

func makeGetLandingsEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetAllLandingsReq)

		ls, err := s.GetAllLandings(
			ctx,
			req.Start,
			req.Count,
			req.Tag,
			req.Name,
		)
		if err != nil {
			return nil, err
		}

		return &GetAllLandingResp{
			Landings: ls,
		}, nil
	}
}