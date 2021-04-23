package estimates

import (
	"context"
	"time"

	"github.com/ITLab-Projects/pkg/models/estimate"
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

	e := estimate.Estimate{}

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

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	return er.DeleteMany(
		ctx,
		bson.M{"milestone_id": bson.M{"$nin": ids}},
		func(dr *mongo.DeleteResult) error {
			return nil
		},
		options.Delete(),
	)
}

func (er *EstimateRepository) save(v interface{}) error {
	estimate, _ := v.(model.Estimate)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": estimate.MilestoneID}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := er.estimateColletion.ReplaceOne(ctx, filter, estimate, opts)
	if err != nil {
		return err
	}

	return nil
}

func (er *EstimateRepository) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
