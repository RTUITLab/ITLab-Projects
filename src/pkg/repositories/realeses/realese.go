package realeses

import (
	"context"
	"time"

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
	rls, _ := v.([]model.RealeseInRepo)

	var ids []uint64
	for _, r := range rls {
		ids = append(ids, r.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (r *RealeseRepo) Save(v interface{}) error {
	err := r.Saver.Save(v)

	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseRepo) save(v interface{}) error {
	real, _ := v.(model.RealeseInRepo)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": real.ID}
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := r.realeseCollection.ReplaceOne(ctx, filter, real, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseRepo) saveAll(rls []model.RealeseInRepo) error {
	for _, rl := range rls {
		if err := r.save(rl); err != nil {
			return err
		}
	}

	return nil
}


func (r *RealeseRepo) SaveAndDeletedUnfind(ctx context.Context, rls interface{}) error {
	if err := r.Saver.SaveAndDeletedUnfind(ctx, rls); err != nil {
		return err
	}

	return nil
}


