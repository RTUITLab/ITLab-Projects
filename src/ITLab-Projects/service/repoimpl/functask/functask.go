package functask

import (
	a "github.com/ITLab-Projects/service/repoimpl/assetsformilestone"
	"context"

	model "github.com/ITLab-Projects/pkg/models/functask"

	"github.com/ITLab-Projects/pkg/repositories/functasks"
)


type FuncTaskRepositoryImp struct {
	FuncTask	functasks.FuncTaskRepositorier
	m			a.MilestoneAssets
}

func New(
	FuncTask functasks.FuncTaskRepositorier,
) *FuncTaskRepositoryImp {
	return &FuncTaskRepositoryImp{
		FuncTask: FuncTask,
		m: a.New(FuncTask),
	}
}

func (f *FuncTaskRepositoryImp) GetFuncTaskByMilestoneID(
	ctx 		context.Context,
	MilestoneID	uint64,
) (*model.FuncTaskFile, error) {
	var task model.FuncTaskFile

	if err := f.m.GetByMilestoneID(
		ctx,
		MilestoneID,
		&task,
	); err != nil {
		return nil, err
	}

	return &task, nil
}

func (f *FuncTaskRepositoryImp) SaveFuncTask(
	ctx context.Context,
	// Can be slice or single value
	task interface{},
) error {
	return f.m.Save(
		ctx,
		task,
	)
}

func (f *FuncTaskRepositoryImp) DeleteOneFuncTaskByMilestoneID(
	ctx 		context.Context,
	MilestoneID uint64,
) error {
	return f.m.DeleteOneByMilestoneID(
		ctx,
		MilestoneID,
	)
}

func (f *FuncTaskRepositoryImp) DeleteManyFuncTasksByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) error {
	return f.m.DeleteManyByMilestoneID(
		ctx,
		MilestonesID,
	)
}

func (f *FuncTaskRepositoryImp) GetFuncTasksByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) ([]*model.FuncTaskFile, error) {
	var fts []*model.FuncTaskFile

	if err := f.m.GetManyByMilestonesID(
		ctx,
		MilestonesID,
		&fts,
	); err != nil {
		return nil, err
	}

	return fts, nil
}