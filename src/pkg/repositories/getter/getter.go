package getter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Getter interface {
	GetAllFiltered(
		ctx context.Context, 
		filter interface{}, 
		f func(*mongo.Cursor) error, 
		opts ...*options.FindOptions,
	) error
	GetOne(
		ctx context.Context, 
		filter interface{}, 
		f func(*mongo.SingleResult) error, 
		opts ...*options.FindOneOptions,
	) error
	GetAll(reps interface{}) error
}