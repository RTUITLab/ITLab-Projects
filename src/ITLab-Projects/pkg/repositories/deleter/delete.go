package deleter

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type Delete struct {
	collection *mongo.Collection
}

func New(c *mongo.Collection) Deleter {
	return &Delete{
		collection: c,
	}
}

func (d *Delete) DeleteMany(
	ctx context.Context,
	filter interface{},
	// if nil would'nt call
	f func(*mongo.DeleteResult) error,
	opts ...*options.DeleteOptions,
) error {
	res, err := d.collection.DeleteMany(
		ctx,
		filter,
		opts...
	)

	if err != nil {
		return err
	}

	if f != nil {
		return f(res)
	}

	return nil
}

func (d *Delete) DeleteOne(
	ctx context.Context,
	filter interface{},
	// if nil would'nt call
	f func(*mongo.DeleteResult) error,
	opts ...*options.DeleteOptions,
) error {
	res, err := d.collection.DeleteOne(
		ctx,
		filter,
		opts...
	)

	if err != nil {
		return err
	}

	if f != nil {
		return f(res)
	}

	return nil
}