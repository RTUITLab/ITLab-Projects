package deleter

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
)

type Deleter interface {
	Delete(
		ctx context.Context,
		filter interface{},
		f func(*mongo.DeleteResult) error,
		opts ...*options.DeleteOptions,
	) error
}