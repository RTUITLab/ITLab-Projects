package issues

import (
	"context"

	"github.com/ITLab-Projects/pkg/chunkresp/linkbuilder"
	"github.com/ITLab-Projects/pkg/conextvalue/chunck"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetIssues	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetIssues: makeGetIssuesEndpoint(s),
	}
}

func makeGetIssuesEndpoint(
	s Service,
) endpoint.Endpoint {
	return func(
		ctx context.Context, 
		request interface{},
	) (response interface{}, err error) {
		req := request.(*GetIssuesReq)

		c := &chunkresp.ChunkResp{}
		ctx = chunck.New(
			ctx,
			c,
		)

		is, err := s.GetIssues(
			ctx,
			req.Start,
			req.Count,
			req.Name,
			req.Tag,
		)
		if err != nil {
			return nil, err
		}
		
		resp := &GetIssuesResp{}

		resp.Issues = is
		resp.ChunkResp = c

		resp.ChunkResp.WriteLinks(
			linkbuilder.New(
				"start",
				"count",
			).Build(
				resp.ChunkResp,
				req.HttpURL,
			),
		)

		return resp, nil
	}
}