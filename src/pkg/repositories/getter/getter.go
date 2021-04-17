package getter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Getter interface {
	GetOner
	GetAller
	GetAllerFiltered
}

type GetOner interface {
	GetOne(
		ctx context.Context, 
		filter interface{}, 
		f func(*mongo.SingleResult) error, 
		opts ...*options.FindOneOptions,
	) error
}

type GetAller interface {
	GetAll(reps interface{}) error
}

type GetAllerFiltered interface {
	GetAllFiltered(
		ctx context.Context, 
		filter interface{}, 
		f func(*mongo.Cursor) error, 
		opts ...*options.FindOptions,
	) error
}