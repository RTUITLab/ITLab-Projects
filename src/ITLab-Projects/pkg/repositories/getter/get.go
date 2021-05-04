package getter

import (
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
)

type Get struct {
	collection 	*mongo.Collection
	checker		typechecker.TypeChecker
}

func New(c *mongo.Collection, typeChecker typechecker.TypeChecker) Getter {
	return &Get{
		collection: c,
		checker: typeChecker,
	}
}

func (g *Get) GetAllFiltered(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.Cursor) error, 
	opts ...*options.FindOptions,
) error {
	cur, err := g.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	if err := f(cur); err != nil {
		return err
	}

	return nil
}

func (g *Get) GetOne(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.SingleResult) error, 
	opts ...*options.FindOneOptions,
) error {
	single := g.collection.FindOne(ctx, filter, opts...)
	if single.Err() != nil {
		return single.Err()
	}

	if err := f(single); err != nil {
		return err
	}
	return nil
}

