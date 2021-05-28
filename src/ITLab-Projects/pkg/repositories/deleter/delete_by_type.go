package deleter

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"github.com/Kamva/mgm"
)

type DeleteByType struct {
	_type mgm.Model
}

func NewDeleteByType(
	Type mgm.Model,
) *DeleteByType {
	return &DeleteByType{
		_type: Type,
	}
}

func (d *DeleteByType) DeleteMany(
	ctx context.Context,
	filter interface{},
	f func(*mongo.DeleteResult) error,
	opts ...*options.DeleteOptions,
) error {
	res, err := mgm.Coll(d._type).DeleteMany(
		ctx,
		filter,
		opts...
	)
	if err != nil {
		return err
	}
	
	if f != nil {
		if err := f(res); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeleteByType) DeleteOne(
	ctx context.Context,
	filter interface{},
	f func(*mongo.DeleteResult) error,
	opts ...*options.DeleteOptions,
) error {
	res, err := mgm.Coll(d._type).DeleteOne(
		ctx,
		filter,
		opts...
	)
	if err != nil {
		return err
	}

	if f != nil {
		if err := f(res); err != nil {
			return err
		}
	}

	return nil
}