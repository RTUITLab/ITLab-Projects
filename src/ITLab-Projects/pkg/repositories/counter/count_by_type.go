package counter

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CountByType struct {
	_type mgm.Model
	count int64
}

func NewCountByType(
	Type mgm.Model,
) *CountByType {
	c := &CountByType{
		_type: Type,
	}

	return c
}

func (c *CountByType) Count() int64 {
	return c.count
}

func (c *CountByType) UpdateCount() (int64, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	count, err := mgm.Coll(c._type).CountDocuments(
		ctx,
		bson.M{},
		options.Count(),
	)
	if err != nil {
		return 0, nil
	}

	c.count = count

	return count, nil
}

func (c *CountByType) CountByFilter(
	filter interface{},
) (int64, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	count, err := mgm.Coll(c._type).CountDocuments(
		ctx,
		filter,
		options.Count(),
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}