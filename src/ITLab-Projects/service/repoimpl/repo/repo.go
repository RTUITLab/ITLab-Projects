package repo

import (
	"github.com/ITLab-Projects/service/repoimpl/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/ITLab-Projects/pkg/models/repo"
	"context"

	"github.com/ITLab-Projects/pkg/repositories/repos"
)

type RepoRepositoryImp struct {
	Repo repos.ReposRepositorier	
}

func New(
	Repo repos.ReposRepositorier,
) *RepoRepositoryImp {
	return &RepoRepositoryImp{
		Repo: Repo,
	}
}

func (r *RepoRepositoryImp) SaveReposAndSetDeletedUnfind(
	ctx context.Context,
	repos interface{},
) error {
	return utils.SaveAndSetDeletedUnfind(
		ctx,
		r.Repo,
		repos,
	)
}

func (r *RepoRepositoryImp) GetFiltrSortRepos(
	ctx 	context.Context,
	filter 	interface{},
	sort 	interface{},
) ([]*repo.Repo, error) {
	return r.GetRepos(
		ctx,
		filter,
		options.Find().
			SetSort(sort),
	)
}

func (r *RepoRepositoryImp) GetReposAndScanTo(
	ctx context.Context,
	filter interface{},
	value interface{},
	options ...*options.FindOptions,
) error {
	return r.Repo.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			if c.RemainingBatchLength() == 0 {
				return mongo.ErrNoDocuments
			}
			return c.All(
				ctx,
				value,
			)
		},
		options...
	)
}

func (r *RepoRepositoryImp) GetRepos(
	ctx 	context.Context,
	filter 	interface{},
	options ...*options.FindOptions,
) ([]*repo.Repo, error) {
	var repos []*repo.Repo

	if err := r.GetReposAndScanTo(
		ctx,
		filter,
		&repos,
		options...
	); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *RepoRepositoryImp) GetFilteredRepos(
	ctx 	context.Context,
	filter 	interface{},
) ([]*repo.Repo, error) {
	return r.GetRepos(
		ctx,
		filter,
		options.Find(),
	)
}

func (r *RepoRepositoryImp) GetFiltrSortFromToRepos(
	ctx 	context.Context,
	filter 	interface{},
	sort 	interface{},
	start 	int64,
	count 	int64,
) ([]*repo.Repo, error) {
	return r.GetRepos(
		ctx,
		filter,
		options.Find().
			SetSort(sort).
			SetSkip(start).
			SetLimit(count),
	)
}

func (r *RepoRepositoryImp) GetByID(
	ctx context.Context,
	ID uint64,
) (*repo.Repo, error) {
	var rep repo.Repo

	if err := r.Repo.GetOne(
		ctx,
		bson.M{"id": ID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&rep)
		},
		options.FindOne(),
	); err != nil {
		return nil, err
	}

	return &rep, nil
}

func (r *RepoRepositoryImp) DeleteByID(
	ctx context.Context,
	ID uint64,
) error {
	if err := r.Repo.DeleteOne(
		ctx,
		bson.M{"id": ID},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}

			return nil
		},
		options.Delete(),
	); err != nil {
		return err
	}

	return nil
}