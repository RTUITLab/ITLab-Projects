package repos

import (
	"context"
	"time"

	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"

	"github.com/ITLab-Projects/pkg/models/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReposRepository struct {
	repoCollection 		*mongo.Collection
	CountOfDocuments	int64
	counter.Counter
	getter.Getter
	saver.Saver
}


func New(repoCollection *mongo.Collection) ReposRepositorier {
	rr := &ReposRepository{
		repoCollection: repoCollection,
		Counter: counter.New(repoCollection),
	}

	Type := repo.Repo{}

	rr.Getter = getter.New(
		repoCollection,
		typechecker.NewSingleByInterface(Type),
	)

	rr.Saver = saver.New(
		repoCollection,
		Type,
		rr.save,
	)

	return rr
}

func (r *ReposRepository) Save(repos interface{}) error {
	if err := r.Saver.Save(repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *ReposRepository) saveAll(repos []repo.Repo) error {
	for _, rep := range repos {
		if err := r.save(rep); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReposRepository) save(v interface{}) error {
	repo, _ := v.(repo.Repo)
	
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": repo.ID}
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := r.repoCollection.ReplaceOne(ctx, filter, repo, opts)
	if err != nil {
		return err
	}

	return nil
}

