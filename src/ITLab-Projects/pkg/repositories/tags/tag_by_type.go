package tags

import (
	"context"

	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/agregate"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagsByType struct {
	model mgm.Model
	saver.SaverWithDelUpdate
	getter.Getter
	deleter.Deleter
	agregate.Agregater
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

	tt.Agregater = agregate.NewByType(&t)

	return tt
}

func (tg *TagsByType) save(ctx context.Context, v interface{}) error {
	tag := getPointer(v)

	_, err := mgm.Coll(tg.model).ReplaceOne(
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

func (tg *TagsByType) buildFilter(v interface{}) interface{} {
	var ids []uint64

	if tgs, ok := v.([]*model.Tag); ok {
		for _, t := range tgs {
			ids = append(ids, t.RepoID)
		}
	} else if tgs, ok := v.([]model.Tag); ok {
		for _, t := range tgs {
			ids = append(ids, t.RepoID)
		}
	}


	return bson.M{"repo_id": bson.M{"$nin": ids}}
}