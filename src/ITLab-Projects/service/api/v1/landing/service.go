package landing

import (
	"github.com/ITLab-Projects/pkg/models/landing"
	"context"
)

type Service interface {
	GetAllLandings(
		ctx				context.Context,
		start,	count	int64,
		tag,	name	string,
	) ([]*landing.LandingCompact, error)

	GetLanding(
		ctx		context.Context,
		ID		uint64,
	) (*landing.Landing, error)
}