package projects

import (
	"context"
	"fmt"

	"github.com/ITLab-Projects/pkg/conextvalue/chunck"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetProjects	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{

	}
}

func makeGetProjectsEndpoint(
	s Service,
) endpoint.Endpoint {
	return func(
		ctx context.Context,
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetProjectsReq)

		c := &chunkresp.ChunkResp{}
		ctx = chunck.New(
			ctx,
			c,
		)

		projs, err := s.GetProjects(
			ctx,
			req.Start,
			req.Count,
			req.Name,
			req.Tag,
		)
		if err != nil {
			return nil, err
		}

		resp := &GetProjectsResp{}

		resp.Projects = projs
		resp.ChunkResp = c

		resp.ChunkResp.WriteLinks(
			chunkresp.NewLink().
							AddSelf(
								fmt.Sprintf(
									"%s?%s",
									req.HttpURL.Path,
									req.HttpURL.Query().Encode(),
								),
							),
		)


		return resp, nil
	}
}