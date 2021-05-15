package projects

import (
	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"context"
	"github.com/go-kit/kit/log"
)

type service struct {
	repository 	Repository
	logger 		log.Logger
	requester	githubreq.Requester
}


func (s *service) GetProject(
	ctx context.Context,
	ID uint64,
) (*repoasproj.RepoAsProj, error) {
	panic("not implemented") // TODO: Implement
}

func (s *service) GetProjects(ctx context.Context, 
	start int64, 
	count int64,
) ([]*repoasproj.RepoAsProjCompact, error) {
	panic("not implemented") // TODO: Implement
}

func (s *service) UpdateProjects(
	ctx context.Context,
) error {
	panic("not implemented") // TODO: Implement
}

func (s *service) DeleteProject(
	ctx context.Context, 
	ID uint64,
) error {
	panic("not implemented") // TODO: Implement
}

