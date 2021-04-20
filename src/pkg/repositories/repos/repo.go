package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
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
	saver.SaverWithDelete
	deleter.Deleter
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

	rr.SaverWithDelete = saver.NewSaverWithDelete(
		repoCollection,
		Type,
		rr.save,
		rr.buildFilter,
	)

	rr.Deleter = deleter.New(
		repoCollection,
	)

	return rr
}

func (r *ReposRepository) buildFilter(v interface{}) interface{} {
	repos, _ := v.([]repo.Repo)

	var ids []uint64

	for _, rep := range repos {
		ids = append(ids, rep.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (r *ReposRepository) Save(repos interface{}) error {
	if err := r.SaverWithDelete.Save(repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *ReposRepository) SaveAndDeletedUnfind(ctx context.Context, repos interface{}) error {
	if err := r.SaverWithDelete.SaveAndDeletedUnfind(ctx, repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *ReposRepository) deleteExpectNew(v interface{}) error {
	repos, ok := v.([]repo.Repo)
	if !ok {
		return fmt.Errorf("Unexpected type %T Expected %T", v, []repo.Repo{})
	}

	var ids []uint64
	for _, rep := range repos {
		ids = append(ids, rep.ID)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	opt := options.Delete()
	filter := bson.M{"id": bson.M{"$nin": ids}}
	if err := r.Delete(
		ctx,
		filter,
		nil,
		opt,
	); err != nil {
		return err
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

