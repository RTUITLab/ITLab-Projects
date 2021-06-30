package projects

import (
	"context"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/repoasproj"
)

type Service interface {
	GetProject(
		ctx context.Context,
		ID	uint64,
	) (*repoasproj.RepoAsProjPointer, error)

	GetProjects(
		ctx 			context.Context,
		Query			GetProjectsQuery,
	) ([]*repoasproj.RepoAsProjCompactPointers, error)

	DeleteProject(
		ctx context.Context,
		ID	uint64,
		r *http.Request,
	) error
}