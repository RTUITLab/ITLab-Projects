package realeses

import (
	"context"

	model "github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RealeseRepo struct {
	realeseCollection *mongo.Collection
	getter.GetOner
	Saver saver.SaverWithDelete
	deleter.Deleter
}

func New(collection *mongo.Collection) Realeser {
	r := &RealeseRepo{
		realeseCollection: collection,
	}

	m := model.RealeseInRepo{}

	r.GetOner = getter.New(
		r.realeseCollection, 
		typechecker.NewSingleByInterface(m),
	)

	r.Saver = saver.NewSaverWithDelete(
		collection,
		m,
		r.save,
		r.buildFilter,
	)

	r.Deleter = deleter.New(
		collection,	
	)

	return r
}

func (r *RealeseRepo) buildFilter(v interface{}) interface{} {
	rls, _ := v.([]*model.RealeseInRepo)

	var ids []uint64
	for _, r := range rls {
		ids = append(ids, r.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (r *RealeseRepo) Save(ctx context.Context, v interface{}) error {
	err := r.Saver.Save(ctx, v)

	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseRepo) save(ctx context.Context, v interface{}) error {
	real := getPointer(v)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"repoid": real.ID}

	
	_, err := r.realeseCollection.ReplaceOne(ctx, filter, real, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseRepo) SaveAndDeletedUnfind(ctx context.Context, rls interface{}) error {
	if err := r.Saver.SaveAndDeletedUnfind(ctx, rls); err != nil {
		return err
	}

	return nil
}

func getPointer(v interface{}) *model.RealeseInRepo {
	r, ok := v.(*model.RealeseInRepo)
	if !ok {
		_r, _ := v.(model.RealeseInRepo)
		r = &_r
	}

	return r
}

