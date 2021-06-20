package updater

import (
	"context"
)

type Service interface {
	UpdateProjects(
		ctx context.Context,
	) error
}