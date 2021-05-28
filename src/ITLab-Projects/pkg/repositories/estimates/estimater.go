package estimates

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/milestoneassets"
)

type EstimateRepositorier interface {
	milestoneassets.AssetsRepositorier
	Delete(MilestoneID uint64) error
	DeleteEstimatesNotIn([]milestone.MilestoneInRepo) error
}