package milestones

import (
	"context"
	"time"

	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"github.com/sirupsen/logrus"

	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MilestoneRepository struct {
	milestoneCollection *mongo.Collection
	counter.Counter
	getter.Getter
	Saver saver.SaverWithDelUpdate
	deleter.Deleter
}

func New(collection *mongo.Collection) Milestoner {
	mr :=  &MilestoneRepository {
		milestoneCollection: collection,
		Counter: counter.New(collection),
	}

	m := model.MilestoneInRepo{}

	mr.Getter = getter.New(
		collection,
		typechecker.NewSingleByInterface(m),
	)

	mr.Saver = saver.NewSaverWithDelUpdate(
		collection,
		m,
		mr.save,
		mr.buildFilter,
	)

	mr.Deleter = deleter.New(
		collection,
	)

	return mr
}

func (m *MilestoneRepository) buildFilter(v interface{}) interface{} {
	ms, ok := v.([]model.MilestoneInRepo)
	if !ok {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/milestones",
				"func": "buildfilter",
			},
		).Panic()
	}

	var ids []uint64
	for _, m := range ms {
		ids = append(ids, m.Milestone.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (m *MilestoneRepository) Save(milestone interface{}) error {
	err := m.Saver.Save(milestone)
	if err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (m *MilestoneRepository) save(v interface{}) error {
	milestone, _ := v.(model.MilestoneInRepo)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": milestone.Milestone.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
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

func (m *MilestoneRepository) SaveAndDeletedUnfind(
	ctx context.Context,
	milestones interface{},
) error {
	if err := m.Saver.SaveAndDeletedUnfind(
		ctx,
		milestones,
	); err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (m *MilestoneRepository) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := m.Saver.SaveAndUpdatenUnfind(ctx, v, updateFilter); err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}
