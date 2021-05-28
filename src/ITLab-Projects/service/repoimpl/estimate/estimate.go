package estimate

import (
	"context"
	model "github.com/ITLab-Projects/pkg/models/estimate"
	a "github.com/ITLab-Projects/service/repoimpl/assetsformilestone"

	"github.com/ITLab-Projects/pkg/repositories/estimates"
)

type EstimateRepositoryImp struct {
	Estimate 	estimates.EstimateRepositorier
	m			a.MilestoneAssets
}

func New(
	Estimate estimates.EstimateRepositorier,
) *EstimateRepositoryImp {
	return &EstimateRepositoryImp{
		Estimate: Estimate,
		m: a.New(Estimate),
	}
}

func (e *EstimateRepositoryImp) GetEstimateByMilestoneID(
	ctx 		context.Context,
	MilestoneID	uint64,
) (*model.EstimateFile, error) {
	var estimate model.EstimateFile

	if err := e.m.GetByMilestoneID(
		ctx,
		MilestoneID,
		&estimate,
	); err != nil {
		return nil, err
	}

	return &estimate, nil
}

func (e *EstimateRepositoryImp) SaveEstimate(
	ctx context.Context,
	// Can be slice or single value
	estimate interface{},
) error {
	return e.m.Save(
		ctx,
		estimate,
	)
}

func (e *EstimateRepositoryImp) DeleteOneEstimateByMilestoneID(
	ctx 		context.Context,
	MilestoneID uint64,
) error {
	return e.m.DeleteOneByMilestoneID(
		ctx,
		MilestoneID,
	)
}

func (e *EstimateRepositoryImp) DeleteManyEstimatesByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) error {
	return e.m.DeleteManyByMilestoneID(
		ctx,
		MilestonesID,
	)
}

func (e *EstimateRepositoryImp) GetEstimatesByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) ([]*model.EstimateFile, error) {
	var es []*model.EstimateFile

	if err := e.m.GetManyByMilestonesID(
		ctx,
		MilestonesID,
		&es,
	); err != nil {
		return nil, err
	}

	return es, nil
}