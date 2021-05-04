package functask

import (
	"github.com/ITLab-Projects/pkg/models/milestonefile"
)

type FuncTask struct {
	MilestoneID	uint64	`json:"milestone_id" bson:"milestone_id"`
	FuncTaskURL	string	`json:"func_task_url" bson:"func_task_url"`
}

type FuncTaskFile struct {
	milestonefile.MilestoneFile	`bson:",inline"`
}

func (f *FuncTaskFile) CollectionName() string {
	return "functask"
}