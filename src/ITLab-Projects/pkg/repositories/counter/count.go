package counter

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Count struct {
	collection *mongo.Collection
	CountOfDocuments int64
}

func New(collection *mongo.Collection) *Count {
	return &Count{
		collection: collection,
	}
}

func (c *Count) Count() int64 {
	return c.CountOfDocuments
}

func (c *Count) UpdateCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	count, err := c.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	c.CountOfDocuments = count

	return count, nil
}

