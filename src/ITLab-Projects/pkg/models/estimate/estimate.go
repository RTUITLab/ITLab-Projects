package estimate

import (
	"github.com/ITLab-Projects/pkg/models/milestonefile"
)

type Estimate struct {
	MilestoneID	uint64	`json:"milestone_id" bson:"milestone_id"`
	EstimateURL	string	`json:"estimate_url" bson:"estimate_url"`
}

type EstimateFile struct {
	milestonefile.MilestoneFile `bson:",inline"`
}

func (e *EstimateFile) CollectionName() string {
	return "estimate"
}