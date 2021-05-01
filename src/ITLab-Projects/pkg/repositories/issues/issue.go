package issues

import (
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type IssueRepository struct {
	issueCollection *mongo.Collection
	counter.Counter
	Saver saver.SaverWithDelUpdate
	deleter.Deleter
	getter.Getter
}

func New(
	collection *mongo.Collection,
) *IssueRepository  {
	i := &IssueRepository {
		issueCollection: collection,
		Counter: counter.New(collection),
	}

	i.Deleter = deleter.New(
		collection,
	)

	issue := model.IssuesWithMilestoneID{}

	i.Saver = saver.NewSaverWithDelUpdate(
		collection,
		issue,
		i.save,
		i.buildFilter,
	)

	i.Getter = getter.New(
		collection,
		typechecker.NewSingleByInterface(issue),
	)

	return i
}

func (i *IssueRepository) buildFilter(v interface{}) interface{} {
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

func (i *IssueRepository) save(v interface{}) error {
	issue, _ := v.(model.IssuesWithMilestoneID)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"id": issue.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := i.issueCollection.ReplaceOne(ctx, filter, issue, opts)
	if err != nil {
		return err
	}

	return nil
}

func (i *IssueRepository) Save(issue interface{}) error {
	err := i.Saver.Save(issue)
	if err != nil {
		return err
	}

	if _, err := i.UpdateCount(); err != nil {
		return err
	}

	return nil
}

func (i *IssueRepository) SaveAndDeletedUnfind(
	ctx context.Context,
	issues interface{},
) error {
	if err := i.Saver.SaveAndDeletedUnfind(
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


func (i *IssueRepository) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := i.Saver.SaveAndUpdatenUnfind(ctx, v, updateFilter); err != nil {
		return err
	}

	if _, err := i.UpdateCount(); err != nil {
		return err
	}

	return nil
}

