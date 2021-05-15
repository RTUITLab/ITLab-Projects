package tag

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/models/tag"
	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/tags"
)

type TagRepositoryImp struct {
	Tag tags.Tager
}

func New(
	Tag tags.Tager,
) *TagRepositoryImp {
	return &TagRepositoryImp{
		Tag: Tag,
	}
}

func (t *TagRepositoryImp) SaveAndDeleteUnfindTags(
	ctx context.Context,
	tgs interface{},
) (error) {
	if err := t.Tag.SaveAndDeletedUnfind(
		ctx,
		tgs,
	); err != nil {
		return err
	}

	return nil
}

func (t *TagRepositoryImp) DeleteTagsByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) (error) {
	return t.Tag.DeleteMany(
		ctx,
		bson.M{"repo_id": RepoID},
		nil,
		options.Delete(),
	)
}

func (t *TagRepositoryImp) GetFilteredTags(
	ctx context.Context,
	filter interface{},
) ([]*model.Tag, error) {
	var tgs []*model.Tag

	if err := t.Tag.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			return c.All(
				ctx,
				&tgs,
			)
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return tgs, nil
}

func (t *TagRepositoryImp) GetFilteredTagsByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) ([]*model.Tag, error) {
	return t.GetFilteredTags(
		ctx,
		bson.M{"repo_id": RepoID},
	)
}

func (t *TagRepositoryImp) GetAllTags(
	ctx context.Context,
) ([]*tag.Tag, error) {
	_tags, err := t.Tag.Distinct(
		ctx,
		"tag",
		bson.M{},
	)
	if err != nil {
		return nil, err
	}

	var tgs []*tag.Tag

	for _, t := range _tags {
		tgs = append(
			tgs, 
			&model.Tag{
				Tag: fmt.Sprint(t),
			},
		)
	}

	return tgs, nil
}