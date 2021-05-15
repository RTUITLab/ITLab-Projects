package issue

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"context"

	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/service/repoimpl/utils"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/repositories/issues"
)

type IssueRepositoryImp struct {
	Issue issues.Issuer
}

func New(
	Issue issues.Issuer,
) *IssueRepositoryImp {
	return &IssueRepositoryImp{
		Issue: Issue,
	}
}

func (i *IssueRepositoryImp) SaveIssuesAndSetDeletedUnfind(
	ctx context.Context,
	is 	[]*model.IssuesWithMilestoneID,
) (error) {
	return utils.SaveAndSetDeletedUnfind(
		ctx,
		i.Issue,
		is,
	)
}
func (i *IssueRepositoryImp) GetIssues(
	ctx 		context.Context,
	filter 		interface{},
	options 	...*options.FindOptions,
) ([]*model.IssuesWithMilestoneID, error) {
	var is []*model.IssuesWithMilestoneID

	if err := i.Issue.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			c.All(
				ctx,
				&is,
			)
			return c.Err()
		},
		options...,
	); err != nil {
		return nil, err
	}

	return is, nil
}

func (i *IssueRepositoryImp) GetFiltrSortIssues(
	ctx 	context.Context,
	filter 	interface{},
	sort 	interface{},
) ([]*model.IssuesWithMilestoneID, error) {
	return i.GetIssues(
		ctx,
		filter,
		options.Find().
			SetSort(sort),
	)
}

func (i *IssueRepositoryImp) GetFilteredIssues(
	ctx context.Context,
	filter interface{},
) ([]*model.IssuesWithMilestoneID, error) {
	return i.GetIssues(
		ctx,
		filter,
		options.Find(),
	)
}

func (i *IssueRepositoryImp) GetFiltrSortedFromToIssues(
	ctx context.Context,
	filter 	interface{},
	sort 	interface{},
	start 	int64,
	count 	int64,
) ([]*model.IssuesWithMilestoneID, error) {
	return i.GetIssues(
		ctx,
		filter,
		options.Find().
			SetSort(sort).
			SetSkip(start).
			SetLimit(count),
	)
}

func (i *IssueRepositoryImp) GetLabelsNameFromOpenIssues(
	ctx context.Context,
) ([]interface{}, error) {
	return i.Issue.Distinct(
		ctx,
		"labels.name",
		bson.M{"state": "open"},
		options.Distinct(),
	)
}

func (i *IssueRepositoryImp) DeleteAllByMilestoneID(
	ctx 		context.Context,
	MilestoneID uint64,
) error {
	return i.Issue.DeleteMany(
		ctx,
		bson.M{"milestone_id": MilestoneID},
		func(dr *mongo.DeleteResult) error {
			if dr.DeletedCount == 0 {
				return mongo.ErrNoDocuments
			}
			return nil
		},
		options.Delete(),
	)
}