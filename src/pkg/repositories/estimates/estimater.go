package estimates

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type EstimateRepositorier interface {
	saver.Saver
	getter.Getter
	deleter.Deleter
	Delete(MilestoneID uint64) error
	DeleteEstimatesNotIn([]milestone.MilestoneInRepo) error
}