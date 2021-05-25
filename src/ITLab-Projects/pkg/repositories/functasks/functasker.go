package functasks

import (
	"github.com/ITLab-Projects/pkg/repositories/milestoneassets"
	"github.com/ITLab-Projects/pkg/models/milestone"
)

type FuncTaskRepositorier interface {
	milestoneassets.AssetsRepositorier
	Delete(uint64) error
	DeleteFuncTasksNotIn([]milestone.MilestoneInRepo) error
}