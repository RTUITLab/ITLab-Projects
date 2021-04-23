package functasks

import (
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

	ft := model.FuncTask{}

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

func (ftr *FuncTaskRepository) save(v interface{}) error {
	functask, _ := v.(model.FuncTask)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": functask.MilestoneID}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := ftr.funcTaskCollection.ReplaceOne(ctx, filter, functask, opts)
	if err != nil {
		return err
	}

	return nil
}

func (ftr *FuncTaskRepository) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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