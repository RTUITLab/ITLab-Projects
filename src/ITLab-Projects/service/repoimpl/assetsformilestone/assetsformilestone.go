package assetsformilestone

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	a "github.com/ITLab-Projects/pkg/repositories/milestoneassets"
)

type MilestoneAssets interface {
	GetByMilestoneID(
		ctx 		context.Context,
		MilestoneID uint64,
		scanTo		interface{},
	) error
	Save(
		ctx 		context.Context,
		value 		interface{},
	) error
	DeleteOneByMilestoneID(
		ctx 		context.Context,
		MilestoneID	uint64,
	) error
	DeleteManyByMilestoneID(
		ctx 			context.Context,
		MilestonesID	[]uint64,
	) error
	GetManyByMilestonesID(
		ctx				context.Context,
		MilestonesID	[]uint64,
		scanTo interface{},
	) error
}

type milestoneAssetsImp struct {
	Asset a.AssetsRepositorier
}

func New(
	Asset a.AssetsRepositorier,
) MilestoneAssets {
	return &milestoneAssetsImp{
		Asset: Asset,
	}
}

func (m *milestoneAssetsImp) GetByMilestoneID(
	ctx 		context.Context,
	MilestoneID uint64,
	scanTo		interface{},
) error {
	return m.Asset.GetOne(
		ctx,
		bson.M{"milestone_id": MilestoneID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(scanTo)
		},
		options.FindOne(),
	)
}

func (m *milestoneAssetsImp) Save(
	ctx 		context.Context,
	value 		interface{},
) error {
	return m.Asset.Save(
		ctx,
		value,
	)
}

func (m *milestoneAssetsImp) DeleteOneByMilestoneID(
	ctx 		context.Context,
	MilestoneID	uint64,
) error {
	return m.Asset.DeleteOne(
		ctx,
		bson.M{"milestone_id": MilestoneID},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}
			return nil
		},
		options.Delete(),
	)
}

func (m *milestoneAssetsImp) DeleteManyByMilestoneID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) error {
	return m.Asset.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": MilestonesID}},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}
			return nil
		},
		options.Delete(),
	)
}

func (m *milestoneAssetsImp) GetManyByMilestonesID(
	ctx				context.Context,
	MilestonesID	[]uint64,
	scanTo 			interface{},
) error {
	return m.Asset.GetAllFiltered(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": MilestonesID}},
		func(c *mongo.Cursor) error {
			if len := c.RemainingBatchLength(); len == 0 {
				return mongo.ErrNoDocuments
			}
			
			return c.All(
				ctx,
				scanTo,
			)
		},
		options.Find(),
	)
}