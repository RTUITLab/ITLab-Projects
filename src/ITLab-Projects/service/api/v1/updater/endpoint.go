package updater

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UpdateProjects	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		UpdateProjects: makeUpdateProjectsEndpoint(s),
	}
}

func makeUpdateProjectsEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		_ = request.(*UpdateProjectsReq)
		err = s.UpdateProjects(ctx)
		if err != nil {
			return nil, err
		}

		return &UpdateProjectsResp{}, nil
	}
}