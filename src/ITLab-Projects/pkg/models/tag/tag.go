package tag

import "github.com/Kamva/mgm"

type Tag struct {
	mgm.DefaultModel	`json:"-" bson:",inline"`
	RepoID 		uint64 	`json:"-" bson:"repo_id"`
	Tag			string	`json:"tag" bson:"tag"`
}

func (tag *Tag) CollectionName() string {
	return "tags"
}