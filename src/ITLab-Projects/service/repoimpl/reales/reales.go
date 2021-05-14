package reales

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/ITLab-Projects/pkg/models/realese"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ITLab-Projects/pkg/repositories/realeses"
)

type RealeseRepositoryImp struct {
	Realese realeses.Realeser
}

func New(
	Realese realeses.Realeser,
) *RealeseRepositoryImp {
	return &RealeseRepositoryImp{
		Realese: Realese,	
	}
}

func (r *RealeseRepositoryImp) SaveRealeses(
	ctx context.Context,
	rs []*model.RealeseInRepo,
) error {
	return r.Realese.Save(
		ctx,
		rs,
	)
}

func (r *RealeseRepositoryImp) GetByRepoID(
	ctx 		context.Context,
	RepoID		uint64,
) (*model.RealeseInRepo, error) {
	relase := &model.RealeseInRepo{}
	if err := r.Realese.GetOne(
		ctx,
		bson.M{"repoid": RepoID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(relase)
		},
		options.FindOne(),
	); err != nil {
		return nil, err
	}

	return relase, nil
}

