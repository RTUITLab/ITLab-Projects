package repos

import (
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"fmt"
	"context"
	"errors"
	"time"

	"github.com/ITLab-Projects/pkg/models/repo"
	wrapper "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReposRepository struct {
	repoCollection 		*mongo.Collection
	CountOfDocuments	int64
	counter.Counter
	getter.Getter
}


func New(repoCollection *mongo.Collection) ReposRepositorier {
	rr := &ReposRepository{
		repoCollection: repoCollection,
		Counter: counter.New(repoCollection),
	}

	rr.Getter = getter.New(
		repoCollection,
		rr.checkType,	
	)

	return rr
}

func (r *ReposRepository) Save(repos interface{}) error {
	var err error = nil
	switch repos.(type) {
	case []repo.Repo:
		err = r.saveAll(repos.([]repo.Repo))
	case repo.Repo:
		err = r.save(repos.(repo.Repo))
	case *repo.Repo:
		err = r.save(*(repos.(*repo.Repo)))
	default:
		err = wrapper.Wrapf(
			errors.New("Uknown type"), 
			"%T Expected %T or %T or %T", 
			repos, []repo.Repo{}, repo.Repo{}, &repo.Repo{}, 
		)
	}
	
	if err != nil {
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

func (r *ReposRepository) save(repo repo.Repo) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": repo.ID}
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := r.repoCollection.ReplaceOne(ctx, filter, repo, opts)
	if err != nil {
		return err
	}

	return nil
}


func (r *ReposRepository) checkType(v interface{}) error {
	var err error = nil

	switch v.(type) {
	case *[]repo.Repo:
		break
	default:
		err = fmt.Errorf("Uknown type: %T Expected: %T", v, &[]repo.Repo{})
	}

	return err
}

