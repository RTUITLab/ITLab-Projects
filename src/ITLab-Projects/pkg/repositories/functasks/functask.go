package functasks

import (
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
	"time"

	model "github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FuncTaskRepository struct {
	funcTaskCollection *mongo.Collection
	saver.Saver
	getter.Getter
	deleter.Deleter
}

func New(
	collection *mongo.Collection,
) FuncTaskRepositorier {
	ftr := &FuncTaskRepository{
		funcTaskCollection: collection,
	}

	ft := model.FuncTaskFile{}

	ftr.Getter = getter.New(
		collection,
		typechecker.NewSingleByInterface(ft),
	)

	ftr.Saver = saver.NewSaver(
		collection,
		ft,
		ftr.save,
	)

	ftr.Deleter = deleter.New(
		collection,
	)


	return ftr
}

func (ftr *FuncTaskRepository) DeleteFuncTasksNotIn(ms []milestone.MilestoneInRepo) error {
	ids := milestone.GetIDS(ms)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	return ftr.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$nin": ids}},
		func(dr *mongo.DeleteResult) error {
			return nil
		},
		options.Delete(),
	)
}

func (ftr *FuncTaskRepository) save(ctx context.Context, v interface{}) error {
	functask := getPointer(v)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": functask.MilestoneID}

	_, err := ftr.funcTaskCollection.ReplaceOne(ctx, filter, functask, opts)
	if err != nil {
		return err
	}

	return nil
}

func getPointer(v interface{}) *model.FuncTaskFile {
	functask, ok := v.(*model.FuncTaskFile)
	if !ok {
		_f, _ := v.(model.FuncTaskFile)
		functask = &_f
	}
	return functask
}

func (ftr *FuncTaskRepository) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := ftr.funcTaskCollection.DeleteOne(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}