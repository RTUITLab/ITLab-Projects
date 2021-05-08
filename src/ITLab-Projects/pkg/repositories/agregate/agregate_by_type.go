package agregate

import (
	"context"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AgregateByType struct {
	_type mgm.Model
}

func NewByType(
	_type mgm.Model,
) *AgregateByType {
	return &AgregateByType{
		_type: _type,
	}
}

func (a *AgregateByType) Agregate(
	ctx context.Context,
	pipeline interface{},
	f func(*mongo.Cursor) error,
	opts... *options.AggregateOptions,
) error {
	cur, err := mgm.Coll(a._type).Aggregate(
		ctx,
		pipeline,
		opts...,
	)
	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	if err := f(cur); err != nil {
		return err
	}

	return nil
}

func (a *AgregateByType) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts... *options.DistinctOptions,
) ([]interface{} ,error) {
	fields, err := mgm.Coll(a._type).Distinct(
		ctx,
		fieldName,
		filter,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return fields, nil
}