package milestone_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/joho/godotenv"
	"github.com/ITLab-Projects/pkg/repositories"
	"testing"

	"github.com/ITLab-Projects/service/repoimpl/milestone"
)


var Repositories *repositories.Repositories
var MilestoneRepository *milestone.MilestoneRepositoryImp
func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}

	Repositories = test.GetTestRepository()

	MilestoneRepository = milestone.New(
		Repositories.Milestone,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init sucess")
}

func TestFunc_SaveMilestoneAndSetdDeletedUnfind(t *testing.T) {
	deleted := []*model.MilestoneInRepo{
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 1, Title: "mock_1",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 2, Title: "mock_2",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 3, Title: "mock_3",
			},
		}},
	}

	not_deleted := []*model.MilestoneInRepo{
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 4, Title: "mock_4",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 5, Title: "mock_5",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 6, Title: "mock_6",
			},
		}},
	}

	all := []*model.MilestoneInRepo{}
	all = append(all, deleted...)
	all = append(all, not_deleted...)

	if err := Repositories.Milestone.Save(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	all_ids := []uint{1,2,3,4,5,6}
	defer Repositories.Milestone.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": all_ids}},
		nil,
		options.Delete(),
	)

	if err := MilestoneRepository.SaveMilestonesAndSetDeletedUnfind(
		context.Background(),
		not_deleted,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	ms, err := MilestoneRepository.GetAllMilestonesInRepo(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, m := range ms {
		if m.Milestone.ID == 1 || m.Milestone.ID == 2 || m.Milestone.ID == 3 {
			if !m.Deleted {
				t.Log("Asserting error: not deleted")
				t.FailNow()
			}
		} else if m.Deleted {
				t.Log("Asserting error: deleted but should'nt")
				t.FailNow()
		}

	}
}

func TestFunc_SaveMilestoneAndSetdDeletedUnfindByValue(t *testing.T) {
	deleted := []model.MilestoneInRepo{
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 1, Title: "mock_1",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 2, Title: "mock_2",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 3, Title: "mock_3",
			},
		}},
	}

	not_deleted := []model.MilestoneInRepo{
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 4, Title: "mock_4",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 5, Title: "mock_5",
			},
		}},
		{RepoID: 12, Milestone: model.Milestone{
			MilestoneFromGH: model.MilestoneFromGH{
				ID: 6, Title: "mock_6",
			},
		}},
	}

	all := []model.MilestoneInRepo{}
	all = append(all, deleted...)
	all = append(all, not_deleted...)

	if err := Repositories.Milestone.Save(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	all_ids := []uint{1,2,3,4,5,6}
	defer Repositories.Milestone.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": all_ids}},
		nil,
		options.Delete(),
	)

	if err := MilestoneRepository.SaveMilestonesAndSetDeletedUnfind(
		context.Background(),
		not_deleted,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	ms, err := MilestoneRepository.GetAllMilestonesInRepo(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, m := range ms {
		if m.Milestone.ID == 1 || m.Milestone.ID == 2 || m.Milestone.ID == 3 {
			if !m.Deleted {
				t.Log("Asserting error: not deleted")
				t.FailNow()
			}
		} else if m.Deleted {
				t.Log("Asserting error: deleted but should'nt")
				t.FailNow()
		}

	}
}

func TestFunc_GetAllByRepoID(t *testing.T) {
	milestones := []*model.MilestoneInRepo{
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 1, Title: "mock_1",					
				},
			},
		},
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 2, Title: "mock_2",					
				},
			},
		},
		{
			RepoID: 13, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 3, Title: "mock_3",					
				},
			},
		},
	}

	if err := MilestoneRepository.SaveMilestonesAndSetDeletedUnfind(
		context.Background(),
		milestones,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer Repositories.Milestone.DeleteMany(
		context.Background(),
		bson.M{},
		func(dr *mongo.DeleteResult) error {
			return nil
		},
		options.Delete(),
	)

	ms, err := MilestoneRepository.GetAllMilestonesInRepo(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, m := range ms {
		if !(m.Milestone.ID == 1 && m.Title == "mock_1" || m.Milestone.ID == 2 && m.Title == "mock_2") {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_DeleteAllByRepoID(t *testing.T) {
	milestones := []*model.MilestoneInRepo{
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 1, Title: "mock_1",					
				},
			},
		},
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 2, Title: "mock_2",					
				},
			},
		},
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 3, Title: "mock_3",					
				},
			},
		},
	}

	if err := MilestoneRepository.SaveMilestonesAndSetDeletedUnfind(
		context.Background(),
		milestones,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := MilestoneRepository.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := MilestoneRepository.GetAllMilestonesInRepo(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_DeleteAllByRepoID_ErrNoDocument(t *testing.T) {
	if err := MilestoneRepository.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetMilestoneByID(t *testing.T) {
	if err := Repositories.Milestone.Save(
		context.Background(),
		model.MilestoneInRepo{
			RepoID: 1,
			Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 1,
					Title: "mock_1",
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer MilestoneRepository.DeleteAllMilestonesByRepoID(
		context.Background(),
		1,
	)

	m, err := MilestoneRepository.GetMilestoneByID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !(m.ID == 1 || m.Title == "mock_1") {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_DeleteMilestone(t *testing.T) {
	if _, err := MilestoneRepository.GetMilestoneByID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetMilestoneAndScanTo(t *testing.T) {
	milestones := []*model.MilestoneInRepo{
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 1, Title: "mock_1",					
				},
			},
		},
		{
			RepoID: 12, Milestone: model.Milestone{
				MilestoneFromGH: model.MilestoneFromGH{
					ID: 2, Title: "mock_2",					
				},
			},
		},
	}

	if err := Repositories.Milestone.Save(
		context.Background(),
		milestones,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer MilestoneRepository.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	)

	var gets []*model.Milestone
	err := MilestoneRepository.GetMilestonesAndScanTo(
		context.Background(),
		bson.M{"repoid": 12},
		&gets,
		options.Find(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, m := range gets {
		if !(m.ID == 1 && m.Title == "mock_1" || m.ID == 2 && m.Title == "mock_2") {
			t.Log("Assert error")
			t.FailNow()
		}
	}	
}

func TestFunc_GetMilestonesAndScanTo_NoDocuments(t *testing.T) {
	err := MilestoneRepository.GetMilestonesAndScanTo(
		context.Background(),
		bson.M{"repoid": 12},
		nil,
		options.Find(),
	)
	if err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}