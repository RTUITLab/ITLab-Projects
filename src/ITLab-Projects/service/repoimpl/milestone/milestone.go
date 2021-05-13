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

func (m *MilestoneRepositoryImp) SaveMilestonesAndSetDeletedUnfind(
	ctx context.Context,
	ms []*model.MilestoneInRepo,
) error {
	return utils.SaveAndSetDeletedUnfind(
		ctx,
		m.Milestone,
		ms,
	)
}

func (m *MilestoneRepositoryImp) GetAllByRepoID(
	ctx 		context.Context,
	RepoID		uint64,
) ([]*model.MilestoneInRepo, error)  {
	var ms []*model.MilestoneInRepo

	if err := m.Milestone.GetAllFiltered(
		ctx,
		bson.M{"repoid": RepoID},
		func(c *mongo.Cursor) error {
			return c.All(
				ctx,
				&ms,
			)
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return ms, nil
}

func (m *MilestoneRepositoryImp) DeleteAllByRepoID(
	ctx 	context.Context,
	RepoID	uint64,
) (error) {
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

// func (m *MilestoneRepositoryImp) CountCompletedForRepo(
// 	ctx 	context.Context,
// 	RepoID	uint64,
// ) (float64,)