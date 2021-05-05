package issues

import (
	"github.com/sirupsen/logrus"
	"context"

	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IssueByType struct {
	saver.SaverWithDelUpdate
	counter.Counter
	deleter.Deleter
	getter.Getter
	model mgm.Model
}

func NewByType(

) *IssueByType {
	it := &IssueByType{}

	i := model.IssuesWithMilestoneID{}
	it.model = &i

	it.Counter = counter.NewCountByType(
		&i,
	)

	it.Deleter = deleter.NewDeleteByType(
		&i,
	)

	it.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		i,
		&i,
		it.save,
		it.buildFilter,
	)

	it.Getter = getter.NewGetByType(
		&i,
	)

	return it
}

func (i *IssueByType) save(ctx context.Context, v interface{}) error {
	issue, _ := v.(model.IssuesWithMilestoneID)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": issue.ID}

	
	_, err := mgm.Coll(i.model).ReplaceOne(ctx, filter, issue, opts)
	if err != nil {
		return err
	}

	return nil
}

func (i *IssueByType) buildFilter(v interface{}) interface{} {
	is, ok := v.([]model.IssuesWithMilestoneID)
	if !ok {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/issues",
				"func": "buildfilter",
			},
		).Panic()
	}

	var ids []uint64
	for _, i := range is {
		ids = append(ids, i.ID)
	}

	return bson.M{"id": bson.M{"$nin": ids}}
}

func (i *IssueByType) Save(ctx context.Context, issue interface{}) error {
	err := i.SaverWithDelUpdate.Save(ctx, issue)
	if err != nil {
		return err
	}

	if _, err := i.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (i *IssueByType) SaveAndDeletedUnfind(
	ctx context.Context,
	issues interface{},
) error {
	if err := i.SaverWithDelUpdate.SaveAndDeletedUnfind(
		ctx,
		issues,	
	); err != nil {
		return err
	}

	if _, err := i.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (i *IssueByType) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := i.SaverWithDelUpdate.SaveAndUpdatenUnfind(ctx, v, updateFilter); err != nil {
		return err
	}

	if _, err := i.UpdateCount(); err != nil {
		return err
	}

	return nil
}