package issue_test

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"os"
	"testing"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/joho/godotenv"
)

var Repositories *repositories.Repositories
var IssueRepository *issue.IssueRepositoryImp

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI_TEST")
	if !find {
		panic("Don't find dburi")
	}

	_r, err := repositories.New(&repositories.Config{
		DBURI: dburi,
	})
	if err != nil {
		panic(err)
	}

	Repositories = _r

	IssueRepository = issue.New(
		_r.Issue,
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init")
}

func TestFunc_SaveIssuesAndSetDeletedUnfind(t *testing.T) {
	deleted := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 1, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 2, Title: "mock_2"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 3, Title: "mock_3"}},
	}

	not_deleted := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 4, Title: "mock_5"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 5, Title: "mock_6"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 6, Title: "mock_7"}},
	}

	all := []*model.IssuesWithMilestoneID{}
	all = append(all, deleted...)
	all = append(all, not_deleted...)

	if err := Repositories.Issue.Save(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	all_ids := []uint{1,2,3,4,5,6}
	defer Repositories.Issue.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": all_ids}},
		nil,
		options.Delete(),
	)

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		not_deleted,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	is, err := IssueRepository.GetIssues(
		context.Background(),
		bson.M{},
		options.Find(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range is {
		if i.ID == 1 || i.ID == 2 || i.ID == 3 {
			if !i.Deleted {
				t.Log("Asserting error: not deleted")
				t.FailNow()
			}
		} else if i.Deleted {
				t.Log("Asserting error: deleted but should'nt")
				t.FailNow()
		}

	}
}