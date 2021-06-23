package landing

import (
	"github.com/ITLab-Projects/pkg/models/landing"
	"context"
)

type Service interface {
	GetAllLandings(
		ctx				context.Context,
		Query			GetAllLandingsQuery,
	) ([]*landing.LandingCompact, error)

	GetLanding(
		ctx		context.Context,
		ID		uint64,
	) (*landing.Landing, error)
}