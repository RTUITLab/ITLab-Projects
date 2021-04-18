package estimate

type Estimate struct {
	MilestoneID	uint	`json:"milestone_id" bson:"milestone_id"`
	EstimateURL	string	`json:"estimate_url" bson:"estimate_url"`
}