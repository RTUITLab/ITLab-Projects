package milestonefile

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MilestoneFile struct {
	MilestoneID	uint64				`json:"milestone_id" bson:"milestone_id"`
	FileID		primitive.ObjectID	`json:"id" bson:"file_id"`
}