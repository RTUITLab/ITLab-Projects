package projects

import (
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"context"
)

type Service interface {
	GetProjects(
		ctx 			context.Context,
		Query			GetProjectsQuery,
	) ([]*repoasproj.RepoAsProjCompactPointers, error)
}