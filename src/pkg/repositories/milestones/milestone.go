package milestones

import (
	"context"
	"time"

	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"

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
	Saver saver.SaverWithDelete
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

	mr.Saver = saver.NewSaverWithDelete(
		collection,
		m,
		mr.save,
		mr.buildFilter,
	)

	return mr
}

func (m *MilestoneRepository) buildFilter(v interface{}) interface{} {
	ms, _ := v.([]model.MilestoneInRepo)

	var ids []uint64
	for _, m := range ms {
		ids = append(ids, m.ID)
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

