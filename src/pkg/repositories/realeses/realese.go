package realeses

import (
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	wrapper "github.com/pkg/errors"
	"errors"
	"fmt"
	model "github.com/ITLab-Projects/pkg/models/realese"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/ITLab-Projects/pkg/repositories/getter"
)

type RealeseRepo struct {
	realeseCollection *mongo.Collection
	getter.GetOner
}

func (r *RealeseRepo) checkType(v interface{}) error {
	var err error = nil

	switch v.(type) {
	case *[]model.RealeseInRepo:
		break
	default:
		err = fmt.Errorf("Uknown type: %T Expected: %T", v, &[]model.RealeseInRepo{})
	}

	return err
}

func New(collection *mongo.Collection) Realeser {
	r := &RealeseRepo{
		realeseCollection: collection,
	}

	r.GetOner = getter.New(r.realeseCollection, r.checkType)

	return r
}

func (r *RealeseRepo) Save(v interface{}) error {
	var err error = nil
	switch v.(type) {
	case []model.RealeseInRepo:
		err = r.saveAll(v.([]model.RealeseInRepo))
	case model.RealeseInRepo:
		err = r.save(v.(model.RealeseInRepo))
	case *model.RealeseInRepo:
		err = r.save(*(v.(*model.RealeseInRepo)))
	default:
		err = wrapper.Wrapf(
			errors.New("Uknown type"), 
			"%T Expected %T or %T or %T", 
			v, []model.RealeseInRepo{}, model.RealeseInRepo{}, &model.RealeseInRepo{}, 
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *RealeseRepo) save(real model.RealeseInRepo) error {
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



