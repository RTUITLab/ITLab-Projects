package estimates

import (
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type EstimateRepositorier interface {
	saver.Saver
	getter.Getter
	deleter.DeleterOne
	Delete(MilestoneID uint64) error
}