package realeses

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	model "github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
)

type RealeseByType struct {
	model mgm.Model
	saver.SaverWithDelUpdate
	getter.Getter
	deleter.Deleter
}

func NewByType(

) *RealeseByType {
	rt := &RealeseByType{}

	r := model.RealeseInRepo{}

	rt.model = &r

	rt.Getter = getter.NewGetByType(
		&r,
	)

	rt.Deleter = deleter.NewDeleteByType(
		&r,
	)

	rt.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		r,
		&r,
		rt.save,
		rt.buildFilter,
	)

	return rt
}

func (r *RealeseByType) save(ctx context.Context, v interface{}) error {
	real := getPointer(v)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"repoid": real.RepoID}
	
	
	_, err := mgm.Coll(r.model).ReplaceOne(ctx, filter, real, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseByType) buildFilter(v interface{}) interface{} {
	var ids []uint64

	if rls, ok := v.([]*model.RealeseInRepo); ok {
		for _, r := range rls {
			ids = append(ids, r.ID)
		}
	} else if rls, ok := v.([]model.RealeseInRepo); ok {
		for _, r := range rls {
			ids = append(ids, r.ID)
		}
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}