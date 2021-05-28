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
		start, 	count 	int64,
		name, 	tag		string,
	) ([]*repoasproj.RepoAsProjCompactPointers, error)

	UpdateProjects(
		ctx context.Context,
	) error

	DeleteProject(
		ctx context.Context,
		ID	uint64,
		r *http.Request,
	) error
}