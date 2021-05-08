package getter

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"github.com/Kamva/mgm"
)

type GetByType struct {
	_type mgm.Model
}

func NewGetByType(
	Type mgm.Model,
) *GetByType {
	return &GetByType{
		_type: Type,
	}
}

func (g *GetByType) GetOne(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.SingleResult) error, 
	opts ...*options.FindOneOptions,
) error {
	res := mgm.Coll(g._type).FindOne(
		ctx,
		filter,
		opts...,
	)
	if res.Err() != nil {
		return res.Err()
	}
	
	if err := f(res); err != nil {
		return err
	}

	return nil
}

func (g *GetByType) GetAllFiltered(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.Cursor) error, 
	opts ...*options.FindOptions,
) error {
	cur, err := mgm.Coll(g._type).Find(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	if f != nil {
		if err := f(cur); err != nil {
			return err
		}
	}

	return nil
}