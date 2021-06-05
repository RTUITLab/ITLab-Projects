package landing

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	model "github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/repositories/landing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LandingRepositoryImp struct {
	Landing		landing.LandingRepositorier
}

func New(
	Landing landing.LandingRepositorier,
) *LandingRepositoryImp {
	return &LandingRepositoryImp{
		Landing: Landing,
	}
}

func (l *LandingRepositoryImp) SaveAndDeleteUnfindLanding(
	ctx context.Context,
	ls interface{},
) error {
	return l.Landing.SaveAndDeletedUnfind(
		ctx,
		ls,
	)
}

func (l *LandingRepositoryImp) DeleteLandingsByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) (error) {
	return l.Landing.DeleteMany(
		ctx,
		bson.M{"repo_id": RepoID},
		nil,
		options.Delete(),
	)
}

func (l *LandingRepositoryImp) GetFilteredLandings(
	ctx context.Context,
	filter interface{},
) ([]*model.Landing, error) {
	var ls []*model.Landing

	if err := l.Landing.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			if c.RemainingBatchLength() == 0 {
				return mongo.ErrNoDocuments
			}

			return c.All(
				ctx,
				&ls,
			)
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return ls, nil
}

func (l *LandingRepositoryImp) GetLandingByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) (*model.Landing, error) {
	landing := model.Landing{}
	err := l.Landing.GetOne(
		ctx,
		bson.M{"repo_id": RepoID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&landing)
		},
	)
	if err != nil {
		return nil, err
	}

	return &landing, nil
}

func (l *LandingRepositoryImp) GetIDsOfReposByLandingTags(
	ctx		context.Context,
	Tags	[]string,
) ([]uint64, error) {
	ls, err := l.GetFilteredLandings(
		ctx,
		bson.M{"tags": bson.M{"$in": Tags}},
	)
	if err != nil {
		return nil, err
	}

	var ids []uint64

	for _, l := range ls {
		ids = append(ids, l.RepoId)
	}

	return ids, nil
}

func (l *LandingRepositoryImp) GetLandingTagsByRepoID(
	ctx		context.Context,
	RepoID	uint64,
) ([]*tag.Tag, error) {
	landing, err := l.GetLandingByRepoID(
		ctx,
		RepoID,
	)
	if err != nil {
		return nil, err
	}

	var tgs []*tag.Tag

	for _, t := range landing.Tags {
		tgs = append(tgs, &tag.Tag{Tag: t})
	}

	return tgs, nil
}

func (l *LandingRepositoryImp) GetAllTags(
	ctx context.Context,
) ([]*tag.Tag, error) {
	_tags, err := l.Landing.Distinct(
		ctx,
		"tags",
		bson.M{},
		options.Distinct(),
	)
	if err != nil {
		return nil, err
	}

	var tags []*tag.Tag

	for _, t := range _tags {
		tags = append(tags, &tag.Tag{Tag: fmt.Sprint(t)})
	}

	return tags, nil
}
