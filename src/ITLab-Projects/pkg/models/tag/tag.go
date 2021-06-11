package tag

type Tag struct {
	RepoID 		uint64 	`json:"-" bson:"repo_id"`
	Tag			string	`json:"tag" bson:"tag"`
}