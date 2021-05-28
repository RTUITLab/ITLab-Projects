package agregate

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Agregater interface {
	Agregate(
		ctx context.Context,
		pipeline interface{},
		f func(*mongo.Cursor) error,
		opts... *options.AggregateOptions,
	) error
	Distinct(
		ctx context.Context,
		fieldName string,
		filter interface{},
		opts... *options.DistinctOptions,
	) ([]interface{} ,error)
}