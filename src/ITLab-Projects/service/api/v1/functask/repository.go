package functask

import (
	"context"

	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/milestone"
)

type Repository interface {
	FuncTaskRepository
	MilestoneRepository
}

type FuncTaskRepository interface {
	GetFuncTaskByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) (*functask.FuncTaskFile, error)

	SaveFuncTask(
		ctx context.Context,
		// Can be slice or single value
		task interface{},
	) error

	DeleteOneFuncTaskByMilestoneID(
		ctx 		context.Context,
		MilestoneID uint64,
	) error
}

type MilestoneRepository interface {
	GetMilestoneByID(
		ctx 		context.Context,
		MilestoneID uint64,
	) (*milestone.Milestone, error)
}