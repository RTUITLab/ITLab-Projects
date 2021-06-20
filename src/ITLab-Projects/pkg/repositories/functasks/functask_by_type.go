package functasks

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	model "github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
)

type FuncTaskByType struct {
	saver.Saver
	getter.Getter
	deleter.Deleter
	model mgm.Model
}

func NewByType(

) *FuncTaskByType {
	ft := &FuncTaskByType{}

	f := model.FuncTaskFile{}
	ft.model = &f

	ft.Saver = saver.NewSaverByType(
		f,
		&f,
		ft.save,
	)

	ft.Getter = getter.NewGetByType(
		&f,
	)

	ft.Deleter = deleter.NewDeleteByType(
		&f,	
	)
	
	return ft
}

func (ftr *FuncTaskByType) save(ctx context.Context, v interface{}) error {
	functask := getPointer(v)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": functask.MilestoneID}

	_, err := mgm.Coll(ftr.model).ReplaceOne(ctx, filter, functask, opts)
	if err != nil {
		return err
	}

	return nil
}

func (ftr *FuncTaskByType) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mgm.Coll(ftr.model).DeleteOne(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ftr *FuncTaskByType) DeleteFuncTasksNotIn(ms []milestone.MilestoneInRepo) error {
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

func getPointer(v interface{}) *model.FuncTaskFile {
	functask, ok := v.(*model.FuncTaskFile)
	if !ok {
		_f, _ := v.(model.FuncTaskFile)
		functask = &_f
	}
	return functask
}