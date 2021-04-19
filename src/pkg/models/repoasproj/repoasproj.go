package repoasproj

import (
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repo"
)

type RepoAsProj struct {
	Repo 			repo.Repo 				`json:"repo"`
	Milestones 		[]milestone.Milestone	`json:"milestones"`
	LastRealese		realese.Realese			`json:"last_realese"`
}