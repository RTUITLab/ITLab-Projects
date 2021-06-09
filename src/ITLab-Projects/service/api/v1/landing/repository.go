package landing

import (
	model "github.com/ITLab-Projects/pkg/models/landing"
	"context"
)


type Repository interface {
	LandingRepository
}

type LandingRepository interface {
	GetFiltrSortLandingCompactFromTo(
		ctx		context.Context,
		filter	interface{},
		sort	interface{},
		start	int64,
		count	int64,
	) ([]*model.LandingCompact, error)

	GetLandingByRepoID(
		ctx 	context.Context,
		RepoID	uint64,
	) (*model.Landing, error)

	GetIDsOfReposByLandingTags(
		ctx		context.Context,
		Tags	[]string,
	) ([]uint64, error)
}