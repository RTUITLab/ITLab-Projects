package projects

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetProject		endpoint.Endpoint
	GetProjects		endpoint.Endpoint
	DeleteProject	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetProject: makeGetProjectEndpoint(s),
		GetProjects: makeGetProjectsEndpoint(s),
		DeleteProject: makeDeleteProjectEndpoint(s),
	}
}

func makeGetProjectEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetProjectReq)
		project, err := s.GetProject(
			ctx,
			req.ID,
		)
		if err != nil {
			return nil, err
		}

		return &GetProjectResp{
			RepoAsProjPointer: project,
		}, nil
	}
}

func makeGetProjectsEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetProjectsReq)
		ps, err := s.GetProjects(
			ctx,
			req.Start,
			req.Count,
			req.Name,
			req.Tag,
		)
		if err != nil {
			return nil, err
		}

		return &GetProjectsResp{
			Projects: ps,
		}, nil
	}
}

func makeDeleteProjectEndpoint(s Service) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*DeleteProjectReq)
		err = s.DeleteProject(
			ctx,
			req.ID,
			req.Req,
		)
		if err != nil {
			return nil, err
		}

		return &DeleteProjectsResp{}, nil
	}
}