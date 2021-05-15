package projects

import (
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"context"
)

type Service interface {
	GetProject(
		ctx context.Context,
		ID	uint64,
	) (*repoasproj.RepoAsProj, error)

	GetProjects(
		ctx 			context.Context,
		start, count 	int64,
	) ([]*repoasproj.RepoAsProjCompact, error)

	UpdateProjects(
		ctx context.Context,
	) error

	DeleteProject(
		ctx context.Context,
		ID	uint64,
	) error
}