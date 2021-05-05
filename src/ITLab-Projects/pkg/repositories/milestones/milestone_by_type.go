package milestones

import (
	"github.com/sirupsen/logrus"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
)

type MilestoneByType struct {
	model mgm.Model
	saver.SaverWithDelUpdate
	counter.Counter
	getter.Getter
	deleter.Deleter
}

func NewByType(

) *MilestoneByType {
	mt := &MilestoneByType{}

	m := model.MilestoneInRepo{}

	mt.model = &m

	mt.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		m,
		&m,
		mt.save,
		mt.buildFilter,
	)

	mt.Counter = counter.NewCountByType(
		&m,
	)

	mt.Getter = getter.NewGetByType(
		&m,	
	)

	mt.Deleter = deleter.NewDeleteByType(
		&m,
	)

	return mt
}

func (m *MilestoneByType) Save(ctx context.Context, milestone interface{}) error {
	if err := m.SaverWithDelUpdate.Save(ctx, milestone); err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (m *MilestoneByType) SaveAndDeletedUnfind(
	ctx context.Context,
	milestones interface{},
) error {
	if err := m.SaverWithDelUpdate.SaveAndDeletedUnfind(ctx,milestones); err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (m *MilestoneByType) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := m.SaverWithDelUpdate.SaveAndUpdatenUnfind(
		ctx,
		v,
		updateFilter,	
	); err != nil {
		return err
	}

	if _, err := m.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (m *MilestoneByType) buildFilter(v interface{}) interface{} {
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

func (m *MilestoneByType) save(ctx context.Context, v interface{}) error {
	milestone, _ := v.(model.MilestoneInRepo)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": milestone.Milestone.ID}
	
	_, err := mgm.Coll(m.model).ReplaceOne(ctx, filter, milestone, opts)
	if err != nil {
		return err
	}

	return nil
}