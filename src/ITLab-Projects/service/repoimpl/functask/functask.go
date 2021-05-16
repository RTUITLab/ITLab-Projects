package functask

import (
	"context"

	model "github.com/ITLab-Projects/pkg/models/functask"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/repositories/functasks"
)

// TODO
type FuncTaskRepositoryImp struct {
	FuncTask	functasks.FuncTaskRepositorier
}

func New(
	FuncTask functasks.FuncTaskRepositorier,
) *FuncTaskRepositoryImp {
	return &FuncTaskRepositoryImp{
		FuncTask: FuncTask,
	}
}

func (f *FuncTaskRepositoryImp) GetFuncTaskByMilestoneID(
	ctx 		context.Context,
	MilestoneID	uint64,
) (*model.FuncTaskFile, error) {
	var task model.FuncTaskFile

	if err := f.FuncTask.GetOne(
		ctx,
		bson.M{"milestone_id": MilestoneID},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&task)
		},
		options.FindOne(),
	); err != nil {
		return nil, err
	}

	return &task, nil
}

func (f *FuncTaskRepositoryImp) SaveFuncTask(
	ctx context.Context,
	// Can be slice or single value
	task interface{},
) error {
	return f.FuncTask.Save(
		ctx,
		task,
	)
}

func (f *FuncTaskRepositoryImp) DeleteOneFuncTaskByMilestoneID(
	ctx 		context.Context,
	MilestoneID uint64,
) error {
	return f.FuncTask.DeleteOne(
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

func (f *FuncTaskRepositoryImp) DeleteManyFuncTasksByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) error {
	return f.FuncTask.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": MilestonesID}},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNilDocument
			}
			return nil
		},
		options.Delete(),
	)
}

func (f *FuncTaskRepositoryImp) GetFuncTasksByMilestonesID(
	ctx 			context.Context,
	MilestonesID	[]uint64,
) ([]*model.FuncTaskFile, error) {
	var fts []*model.FuncTaskFile

	if err := f.FuncTask.GetAllFiltered(
		ctx,
		bson.M{"milestone_id": bson.M{"$in": MilestonesID}},
		func(c *mongo.Cursor) error {
			return c.All(
				ctx,
				&fts,
			)
		},
		options.Find(),
	); err != nil {
		return nil, err
	}

	return fts, nil
}