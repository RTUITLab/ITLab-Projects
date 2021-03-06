package repos

import (
	"context"

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
	counter.Counter
	getter.Getter
	Saver 				saver.SaverWithDelUpdate
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

	rr.Saver = saver.NewSaverWithDelUpdate(
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
	repos, _ := v.([]*repo.Repo)

	var ids []uint64

	for _, rep := range repos {
		ids = append(ids, rep.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (r *ReposRepository) Save(ctx context.Context, repos interface{}) error {
	if err := r.Saver.Save(ctx, repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *ReposRepository) SaveAndDeletedUnfind(ctx context.Context, repos interface{}) error {
	if err := r.Saver.SaveAndDeletedUnfind(ctx, repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *ReposRepository) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := r.Saver.SaveAndUpdatenUnfind(ctx, v, updateFilter); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}


func (r *ReposRepository) save(ctx context.Context, v interface{}) error {
	rep := pointFromInterface(v)
	
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": rep.ID}

	_, err := r.repoCollection.ReplaceOne(ctx, filter, rep, opts)
	if err != nil {
		return err
	}

	return nil
}

func pointFromInterface(v interface{}) *repo.Repo {
	rep, ok := v.(*repo.Repo)
	if !ok {
		_repo, _ := v.(repo.Repo)
		rep = &_repo
	}

	return rep
}