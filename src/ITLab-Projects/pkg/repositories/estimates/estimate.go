package estimates

import (
	"context"
	"time"

	model "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EstimateRepository struct {
	estimateColletion *mongo.Collection
	saver.Saver
	getter.Getter
	deleter.Deleter
}


func New(
	collection *mongo.Collection,
) EstimateRepositorier {
	er := &EstimateRepository{
		estimateColletion: collection,
	}

	e := model.EstimateFile{}

	er.Getter = getter.New(
		collection,
		typechecker.NewSingleByInterface(e),
	)

	er.Saver = saver.NewSaver(
		collection,
		e,
		er.save,
	)

	er.Deleter = deleter.New(
		collection,
	)

	return er
}

func (er *EstimateRepository) DeleteEstimatesNotIn(ms []milestone.MilestoneInRepo) error {
	ids := milestone.GetIDS(ms)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	return er.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$nin": ids}},
		func(dr *mongo.DeleteResult) error {
			return nil
		},
		options.Delete(),
	)
}

func (er *EstimateRepository) save(ctx context.Context, v interface{}) error {
	estimate := getPointer(v)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": estimate.MilestoneID}
	
	_, err := er.estimateColletion.ReplaceOne(ctx, filter, estimate, opts)
	if err != nil {
		return err
	}

	return nil
}

func getPointer(v interface{}) *model.EstimateFile {
	estimate, ok := v.(*model.EstimateFile)
	if !ok {
		_e, _ := v.(model.EstimateFile)
		estimate = &_e
	}
	return estimate
}

func (er *EstimateRepository) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := er.estimateColletion.DeleteOne(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}
