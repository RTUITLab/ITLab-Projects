package tags

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
)

type TagsByType struct {
	model mgm.Model
	saver.SaverWithDelUpdate
	getter.Getter
	deleter.Deleter
}

func NewByType(

) *TagsByType {
	tt := &TagsByType{}

	t := model.Tag{}

	tt.model = &t

	tt.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		t,
		&t,
		tt.save,
		tt.buildFilter,
	)

	tt.Getter = getter.NewGetByType(
		&t,
	)

	tt.Deleter = deleter.NewDeleteByType(
		&t,
	)

	return tt
}

func (tg *TagsByType) save(ctx context.Context, v interface{}) error {
	tag, _ := v.(model.Tag)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"repo_id": tag.RepoID}

	
	_, err := mgm.Coll(tg.model).ReplaceOne(ctx, filter, tag, opts)
	if err != nil {
		return err
	}

	return nil
}

func (tg *TagsByType) buildFilter(v interface{}) interface{} {
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