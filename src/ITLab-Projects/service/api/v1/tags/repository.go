package tags

import (
	"github.com/ITLab-Projects/pkg/models/tag"
	"context"
)

type Repository interface {
	TagRepository
}

type TagRepository interface {
	GetAllTags(
		ctx context.Context,
	) ([]*tag.Tag, error)
}