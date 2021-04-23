package deleter

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
)

type Deleter interface {
	DeleteMany(
		ctx context.Context,
		filter interface{},
		f func(*mongo.DeleteResult) error,
		opts ...*options.DeleteOptions,
	) error
	DeleterOne
}

type DeleterOne interface {
	DeleteOne(
		ctx context.Context,
		filter interface{},
		f func(*mongo.DeleteResult) error,
		opts ...*options.DeleteOptions,
	) error
}