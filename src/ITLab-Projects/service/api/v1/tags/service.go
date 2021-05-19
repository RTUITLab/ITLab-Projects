package tags

import (
	"github.com/ITLab-Projects/pkg/models/tag"
	"context"
)

type Service interface {
	GetAllTags(
		ctx context.Context,
	) ([]*tag.Tag, error)
}