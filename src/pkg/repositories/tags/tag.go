package tags

import (
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	log "github.com/sirupsen/logrus"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagRepository struct {
	tagCollection *mongo.Collection
	saver.SaverWithDelete
	getter.Getter
}

func New(
	collection *mongo.Collection,
) Tager {
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

	return tr
}

func (tg *TagRepository) buildFilter(v interface{}) interface{} {
	tgs, ok := v.([]model.Tag)
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

func (tg *TagRepository) save(v interface{}) error {
	tag, _ := v.(model.Tag)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"repo_id": tag.RepoID}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	_, err := tg.tagCollection.ReplaceOne(ctx, filter, tag, opts)
	if err != nil {
		return err
	}

	return nil
}