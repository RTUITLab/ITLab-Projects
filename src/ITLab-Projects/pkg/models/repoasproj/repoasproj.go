package repoasproj

import (
	"time"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/models/tag"
)

type RepoAsProj struct {
	Repo 			repo.Repo 					`json:"repo"`
	Milestones 		[]milestone.Milestone		`json:"milestones"`
	LastRealese		*realese.Realese			`json:"last_realese"`
	Tags			[]tag.Tag					`json:"tags"`
	Completed		float64						`json:"completed"`
}

type RepoAsProjPointer struct {
	Repo 			*repo.Repo 					`json:"repo"`
	Milestones 		[]*milestone.Milestone		`json:"milestones"`
	LastRealese		*realese.Realese			`json:"last_realese"`
	Tags			[]*tag.Tag					`json:"tags"`
	Completed		float64						`json:"completed"`
}


type RepoAsProjCompact struct {
	Repo 			repo.Repo					`json:"repo"`
	Completed		float64						`json:"completed"`
	Tags			[]tag.Tag					`json:"tags"`
}

type ByCreateDate []RepoAsProjCompact

func (b ByCreateDate) Len() int {
	return len(b)
}

func (b ByCreateDate) Less(i, j int) bool {
	if b[i].Repo.Deleted {
		return false
	} else if b[j].Repo.Deleted {
		return true
	}

	parsedF, _ := time.Parse(time.RFC3339, b[i].Repo.CreatedAt)
	parsedS, _ := time.Parse(time.RFC3339, b[j].Repo.CreatedAt)

	return parsedF.After(parsedS)

}

func (b ByCreateDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type RepoAsProjCompactPointers struct {
	Repo 			*repo.Repo					`json:"repo"`
	Completed		float64						`json:"completed"`
	Tags			[]*tag.Tag					`json:"tags"`
}

type ByCreateDatePointers []*RepoAsProjCompactPointers

func (b ByCreateDatePointers) Len() int {
	return len(b)
}

func (b ByCreateDatePointers) Less(i, j int) bool {
	if b[i].Repo.Deleted {
		return false
	} else if b[j].Repo.Deleted {
		return true
	}

	parsedF, _ := time.Parse(time.RFC3339, b[i].Repo.CreatedAt)
	parsedS, _ := time.Parse(time.RFC3339, b[j].Repo.CreatedAt)

	return parsedF.After(parsedS)

}

func (b ByCreateDatePointers) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
