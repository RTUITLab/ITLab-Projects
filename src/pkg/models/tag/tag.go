package tag

type Tag struct {
	RepoID 		uint64 	`json:"repo_id" bson:"repo_id"`
	Tag			string	`json:"tag" bson:"tag"`
}