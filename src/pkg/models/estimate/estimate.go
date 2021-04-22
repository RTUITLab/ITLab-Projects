package estimate

type Estimate struct {
	MilestoneID	uint64	`json:"milestone_id" bson:"milestone_id"`
	EstimateURL	string	`json:"estimate_url" bson:"estimate_url"`
}