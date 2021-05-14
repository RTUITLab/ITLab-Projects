package repos

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	model "github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
)

type RepoByType struct {
	model mgm.Model
	counter.Counter
	getter.Getter
	saver.SaverWithDelUpdate
	deleter.Deleter
}

func NewByType(

) *RepoByType {
	rt := &RepoByType{}

	r := model.Repo{}
	rt.model = &r

	rt.Counter = counter.NewCountByType(
		&r,
	)

	rt.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		r,
		&r,
		rt.save,
		rt.buildFilter,
	)

	rt.Deleter = deleter.NewDeleteByType(
		&r,
	)

	rt.Getter = getter.NewGetByType(
		&r,
	)
	
	return rt
}

func (r *RepoByType) buildFilter(v interface{}) interface{} {
	repos, _ := v.([]*model.Repo)

	var ids []uint64

	for _, rep := range repos {
		ids = append(ids, rep.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (r *RepoByType) save(ctx context.Context, v interface{}) error {
	rep := pointFromInterface(v)
	
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": rep.ID}
	

	_, err := mgm.Coll(r.model).ReplaceOne(ctx, filter, rep, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoByType) Save(ctx context.Context, repos interface{}) error {
	if err := r.SaverWithDelUpdate.Save(ctx, repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *RepoByType) SaveAndDeletedUnfind(ctx context.Context, repos interface{}) error {
	if err := r.SaverWithDelUpdate.SaveAndDeletedUnfind(ctx, repos); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (r *RepoByType) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := r.SaverWithDelUpdate.SaveAndUpdatenUnfind(ctx, v, updateFilter); err != nil {
		return err
	}

	if _, err := r.UpdateCount(); err != nil {
		return err
	}

	return nil
}