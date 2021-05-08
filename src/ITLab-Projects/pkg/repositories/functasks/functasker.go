package functasks

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type FuncTaskRepositorier interface {
	saver.Saver
	getter.Getter
	deleter.Deleter
	Delete(uint64) error
	DeleteFuncTasksNotIn([]milestone.MilestoneInRepo) error
}