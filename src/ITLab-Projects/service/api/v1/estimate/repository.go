package estimate

import (
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
)

type Repository interface {
	EstimateRepository
	MilestoneRepository
}

type EstimateRepository interface {
	SaveEstimate(
		ctx context.Context,
		// Can be slice or single value
		estimate interface{},
	) error

	DeleteOneEstimateByMilestoneID(
		ctx 		context.Context,
		MilestoneID uint64,
	) error

	GetEstimateByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) (*estimate.EstimateFile, error)
}

type MilestoneRepository interface {
	GetMilestoneByID(
		ctx 		context.Context,
		MilestoneID uint64,
	) (*milestone.Milestone, error)
}