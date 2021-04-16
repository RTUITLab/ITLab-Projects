package milestones

import (
	"fmt"
	"context"
	"errors"
	"time"

	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	wrapper "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MilestoneRepository struct {
	milestoneCollection *mongo.Collection
	counter.Counter
	getter.Getter
}

func New(collection *mongo.Collection) Milestoner {
	mr :=  &MilestoneRepository {
		milestoneCollection: collection,
		Counter: counter.New(collection),
	}

	mr.Getter = getter.New(
		collection,
		mr.checkType,
	)

	return mr
}

func (m *MilestoneRepository) checkType(v interface{}) error {
	var err error = nil

	switch v.(type) {
	case *[]model.MilestoneInRepo:
		break
	default:
		err = fmt.Errorf("Uknown type: %T Expected: %T", v, &[]model.MilestoneInRepo{})
	}

	return err
}

func (m *MilestoneRepository) Save(milestone interface{}) error {
	var err error = nil
	switch milestone.(type) {
	case model.MilestoneInRepo:
		err = m.save(milestone.(model.MilestoneInRepo))
	case *model.MilestoneInRepo:
		err = m.save(*(milestone.(*model.MilestoneInRepo)))
	case []model.MilestoneInRepo:
		err = m.saveAll(milestone.([]model.MilestoneInRepo))
	default:
		err = wrapper.Wrapf(
			errors.New("Uknown type"), 
			"%T Expected %T or %T or %T", 
			milestone, []model.MilestoneInRepo{}, model.MilestoneInRepo{}, &model.MilestoneInRepo{}, 
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *MilestoneRepository) save(milestone model.MilestoneInRepo) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": milestone.ID}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := m.milestoneCollection.ReplaceOne(ctx, filter, milestone, opts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MilestoneRepository) saveAll(milestones []model.MilestoneInRepo) error {
	for _, milestone := range milestones {
		if err := m.save(milestone); err != nil {
			return err
		}
	}

	return nil
}

