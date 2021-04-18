package functask

type FuncTask struct {
	MilestoneID	uint	`json:"milestone_id" bson:"milestone_id"`
	FuncTaskURL	string	`json:"func_task_url" bson:"func_task_url"`
}