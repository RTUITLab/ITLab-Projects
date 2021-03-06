package milestone

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"context"

	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/milestones"
	"github.com/ITLab-Projects/service/repoimpl/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MilestoneRepositoryImp struct {
	Milestone milestones.Milestoner
}

func New(
	Milestone milestones.Milestoner,
) *MilestoneRepositoryImp {
	return &MilestoneRepositoryImp{
		Milestone: Milestone,
	}
}

func (m *MilestoneRepositoryImp) GetMilestonesAndScanTo(
	ctx context.Context,
	filter interface{},
	value interface{},
	opts ...*options.FindOptions,
) error {
	return m.Milestone.GetAllFiltered(
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
		opts...
	)
}

func (m *MilestoneRepositoryImp) SaveMilestonesAndSetDeletedUnfind(
	ctx context.Context,
	ms interface{},
) error {
	return utils.SaveAndSetDeletedUnfind(
		ctx,
		m.Milestone,
		ms,
	)
}

func (m *MilestoneRepositoryImp) GetAllMilestonesInRepo(
	ctx 		context.Context,
	RepoID		uint64,
) ([]*model.MilestoneInRepo, error)  {
	var ms []*model.MilestoneInRepo

	if err := m.Milestone.GetAllFiltered(
		ctx,
		bson.M{"repoid": RepoID},
		func(c *mongo.Cursor) error {
			c.All(
				ctx,
				&ms,
			)

			if len(ms) == 0 {
				return mongo.ErrNoDocuments
			}

			return c.Err()
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return ms, nil
}

func (m *MilestoneRepositoryImp) GetAllMilestonesByRepoID(
	ctx 		context.Context,
	RepoID		uint64,
) ([]*model.Milestone, error)  {
	var ms []*model.Milestone

	if err := m.Milestone.GetAllFiltered(
		ctx,
		bson.M{"repoid": RepoID},
		func(c *mongo.Cursor) error {
			c.All(
				ctx,
				&ms,
			)

			if len(ms) == 0 {
				return mongo.ErrNoDocuments
			}

			return c.Err()
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return ms, nil
}

func (m *MilestoneRepositoryImp) DeleteAllMilestonesByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) error {
	if err := m.Milestone.DeleteMany(
		ctx,
		bson.M{"repoid": RepoID},
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

func (m *MilestoneRepositoryImp) GetMilestoneByID(
	ctx 		context.Context,
	MilestoneID uint64,
) (*model.Milestone, error) {
	ms := &model.Milestone{}

	if err := m.Milestone.GetOne(
		ctx,
		bson.M{"id": MilestoneID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(ms)
		},
		options.FindOne(),
	); err != nil {
		return nil, err
	}

	return ms, nil
}