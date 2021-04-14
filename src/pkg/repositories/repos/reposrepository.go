package repos

import (
	"time"
	"context"

	"github.com/ITLab-Projects/pkg/models/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReposRepository struct {
	repoCollection *mongo.Collection
}


func New(repoCollection *mongo.Collection) ReposRepositorier {
	return &ReposRepository{
		repoCollection: repoCollection,
	}
}

func (r *ReposRepository) Save(repos []repo.Repo) error {
	for _, rep := range repos {
		opts := options.Replace().SetUpsert(true)
		filter := bson.M{"id": rep.ID}
		
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := r.repoCollection.ReplaceOne(ctx, filter, rep, opts)
		// TODO think about should we return err or just print about it
		if err != nil {
			return err
		}
	}

	return nil
}

// func (r *ReposRepository) Get() ([]repo.Repo, error) {
	

// }