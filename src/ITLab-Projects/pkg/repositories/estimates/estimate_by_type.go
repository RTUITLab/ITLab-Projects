package estimates

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"context"
	"time"

	model "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EstimateRepoByType struct {
	saver.Saver
	getter.Getter
	deleter.Deleter
	model mgm.Model
}

func NewByType() *EstimateRepoByType {
	er := &EstimateRepoByType{}

	e := model.EstimateFile{}
	er.model = &e
	er.Saver = saver.NewSaverByType(
		e,
		&e,
		er.save,
	)

	er.Getter = getter.NewGetByType(
		&e,
	)

	er.Deleter = deleter.NewDeleteByType(
		&e,
	)

	return er
}

func (er *EstimateRepoByType) save(ctx context.Context, v interface{}) error {
	estimate := getPointer(v)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"milestone_id": estimate.MilestoneID}
	
	_, err := mgm.Coll(er.model).ReplaceOne(ctx, filter, estimate, opts)
	if err != nil {
		return err
	}

	return nil
}

func (er *EstimateRepoByType) Delete(MilestoneID uint64) error {
	opts := options.Delete()
	filter := bson.M{"milestone_id": MilestoneID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mgm.Coll(er.model).DeleteOne(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}

func (er *EstimateRepoByType) DeleteEstimatesNotIn(ms []milestone.MilestoneInRepo) error {
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