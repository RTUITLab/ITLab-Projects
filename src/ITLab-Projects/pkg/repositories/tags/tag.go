package tags

import (
	"context"

	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagRepository struct {
	tagCollection *mongo.Collection
	saver.SaverWithDelete
	getter.Getter
	deleter.Deleter
}

func New(
	collection *mongo.Collection,
) *TagRepository {
	tr := &TagRepository{
		tagCollection: collection,
	}

	t := model.Tag{}

	tr.SaverWithDelete = saver.NewSaverWithDelete(
		collection,
		t,
		tr.save,
		tr.buildFilter,
	)

	tr.Getter = getter.New(
		collection,
		typechecker.NewSingleByInterface(t),
	)

	tr.Deleter = deleter.New(
		collection,
	)

	return tr
}

func (tg *TagRepository) buildFilter(v interface{}) interface{} {
	tgs, ok := v.([]*model.Tag)
	if !ok {
		log.WithFields(
			log.Fields{
				"package": "repositories/tags",
				"func": "buildFilter",
				"err": "Unable to cast type",
			},
		).Panic()
	}

	var ids []uint64
	for _, t := range tgs {
		ids = append(ids, t.RepoID)
	}

	return bson.M{"repo_id": bson.M{"$nin": ids}}
}

func (tg *TagRepository) save(ctx context.Context, v interface{}) error {
	tag := getPointer(v)
	_, err := tg.tagCollection.ReplaceOne(
		ctx,
		bson.M{"repo_id": tag.RepoID, "tag": tag.Tag},
		tag,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

func getPointer(v interface{}) *model.Tag {
	tag, ok := v.(*model.Tag)
	if !ok {
		_tag, _ := v.(model.Tag)
		tag = &_tag
	}
	return tag
}