package issue_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/ITLab-Projects/pkg/models/label"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	Repositories = test.GetTestRepository()

	IssueRepository = issue.New(
		Repositories.Issue,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
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

func TestFunc_SaveIssuesAndSetDeletedUnfindByValue(t *testing.T) {
	deleted := []model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 1, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 2, Title: "mock_2"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 3, Title: "mock_3"}},
	}

	not_deleted := []model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 4, Title: "mock_5"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 5, Title: "mock_6"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 6, Title: "mock_7"}},
	}

	all := []model.IssuesWithMilestoneID{}
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

func TestFunc_GetFilteredIssues(t *testing.T) {
	all := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 1, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 2, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 3, Title: "mock_3"}},
	}

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	is, err := IssueRepository.GetFilteredIssues(
		context.Background(),
		bson.M{"title": "mock_1"},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range is {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_1") {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetLabalesNameFromOpenIssues(t *testing.T) {
	all := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{
			ID: 1, Title: "mock_1", State: "open", Labels: []label.Label{{CompactLabel: label.CompactLabel{Name: "label_mock_1"}}},},
		},
		{MilestoneID: 1, Issue: model.Issue{
			ID: 2, Title: "mock_2", State: "open", Labels: []label.Label{{CompactLabel: label.CompactLabel{Name: "label_mock_2"}}},},
		},
	}
	
	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	names, err := IssueRepository.GetLabelsNameFromOpenIssues(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
	}

	for _, name := range names {
		if !(fmt.Sprint(name) == "label_mock_1" || fmt.Sprint(name) == "label_mock_2") {
			t.Log("Assert error")
			t.FailNow()
		}
	}
	
}

func TestFunc_DeleteAllByMilesoneID_NoDocument(t *testing.T) {
	if err := IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetFiltrSortIssues(t *testing.T) {
	all := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 1, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 2, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 3, Title: "mock_3"}},
	}

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	is, err := IssueRepository.GetFiltrSortIssues(
		context.Background(),
		bson.M{},
		bson.D{{"id", -1}},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if sorted := sort.SliceIsSorted(
		is,
		func(i, j int) bool {
			return is[i].ID > is[j].ID
		},
	); !sorted {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_GetFiltrSortFromToIssues(t *testing.T) {
	all := []*model.IssuesWithMilestoneID{
		{MilestoneID: 1, Issue: model.Issue{ID: 1, Title: "mock_1"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 2, Title: "mock_2"}},
		{MilestoneID: 1, Issue: model.Issue{ID: 3, Title: "mock_3"}},
	}

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	is, err := IssueRepository.GetFiltrSortedFromToIssues(
		context.Background(),
		bson.M{},
		bson.D{{"id", -1}},
		1,
		2,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(is) != 2 {
		t.Log("Assert error len is not 2")
		t.FailNow()
	}

	is_1, is_2 := is[0], is[1]

	if !(is_1.ID == 2 && is_1.Title == "mock_2" ||is_2.ID == 1 && is_2.Title == "mock_1") {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_GetAllIssuesByMilestoneID_IssueWithMilestoneID(t *testing.T) {
	save_is := []model.IssuesWithMilestoneID{
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 1,
				Title: "mock_1",
			},
		},
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 2,
				Title: "mock_2",
			},
		},
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 3,
				Title: "mock_3",
			},
		},
	}

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		save_is,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	var get_issie_with_id []model.IssuesWithMilestoneID
	if err := IssueRepository.GetIssuesAndScanTo(
		context.Background(),
		bson.M{"milestone_id": 1},
		&get_issie_with_id,
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range get_issie_with_id {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_2" || i.ID == 3 && i.Title == "mock_3") {
			t.Log("Assert error")
			t.FailNow()
		}
	}

	var get_issie_with_id_pointer []*model.IssuesWithMilestoneID
	if err := IssueRepository.GetIssuesAndScanTo(
		context.Background(),
		bson.M{"milestone_id": 1},
		&get_issie_with_id_pointer,
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range get_issie_with_id_pointer {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_2" || i.ID == 3 && i.Title == "mock_3") {
			t.Log("Assert error")
			t.FailNow()
		}
	}

	var get_issue []model.Issue
	if err := IssueRepository.GetIssuesAndScanTo(
		context.Background(),
		bson.M{"milestone_id": 1},
		&get_issue,
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range get_issue {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_2" || i.ID == 3 && i.Title == "mock_3") {
			t.Log("Assert error")
			t.FailNow()
		}
	}

	var get_issue_pointer []*model.Issue
	if err := IssueRepository.GetIssuesAndScanTo(
		context.Background(),
		bson.M{"milestone_id": 1},
		&get_issue,
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range get_issue_pointer {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_2" || i.ID == 3 && i.Title == "mock_3") {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetAllIssuesByMilestoneID(t *testing.T) {
	save_is := []model.IssuesWithMilestoneID{
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 1,
				Title: "mock_1",
			},
		},
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 2,
				Title: "mock_2",
			},
		},
		{
			MilestoneID: 1, 
			Issue: model.Issue{
				ID: 3,
				Title: "mock_3",
			},
		},
	}

	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		save_is,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer IssueRepository.DeleteAllIssuesByMilestoneID(
		context.Background(),
		1,
	)

	is, err := IssueRepository.GetAllIssuesByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range is {
		if !(i.ID == 1 && i.Title == "mock_1" || i.ID == 2 && i.Title == "mock_2" || i.ID == 3 && i.Title == "mock_3") {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_DeleteAllIssuesByMilestoneID(t *testing.T) {
	if err := IssueRepository.SaveIssuesAndSetDeletedUnfind(
		context.Background(),
		[]model.IssuesWithMilestoneID {
			{
				MilestoneID: 1,
				Issue: model.Issue{
					ID: 1,
					Title: "mock_1",
				},
			},
			{
				MilestoneID: 2,
				Issue: model.Issue{
					ID: 2,
					Title: "mock_2",
				},
			},
			{
				MilestoneID: 3,
				Issue: model.Issue{
					ID: 3,
					Title: "mock_3",
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := IssueRepository.DeleteAllIssuesByMilestonesID(
		context.Background(),
		[]uint64{1,2,3},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := IssueRepository.GetAllIssuesByMilestoneID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := IssueRepository.GetAllIssuesByMilestoneID(
		context.Background(),
		2,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := IssueRepository.GetAllIssuesByMilestoneID(
		context.Background(),
		3,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_DeleteAllIssuesByMilesonesID_NoDocuments(t *testing.T) {
	if err := IssueRepository.DeleteAllIssuesByMilestonesID(
		context.Background(),
		[]uint64{1,2,3},
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}