package projects

import (
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"context"
)

type Service interface {
	GetProjects(
		ctx 			context.Context,
		start, 	count 	int64,
		name, 	tag		string,
	) ([]*repoasproj.RepoAsProjCompactPointers, error)
}