package repos

import (
	"github.com/ITLab-Projects/pkg/models/repo"
)

type ReposRepositorier interface {
	Save(repos []repo.Repo) error
	// GetAll() ([]repo.Repo, error)
}