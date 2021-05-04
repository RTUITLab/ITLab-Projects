package milestonefile

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MilestoneFile struct {
	mgm.DefaultModel				`json:"-" bson:",inline"`
	MilestoneID	uint64				`json:"milestone_id" bson:"milestone_id"`
	FileID		primitive.ObjectID	`json:"id" bson:"file_id"`
}