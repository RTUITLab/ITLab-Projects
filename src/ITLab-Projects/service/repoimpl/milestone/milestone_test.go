package milestone_test

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	model "github.com/ITLab-Projects/pkg/models/milestone"
	"os"
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

	MilestoneRepository = milestone.New(
		_r.Milestone,
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

	ms, err := MilestoneRepository.GetAllByRepoID(
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

	ms, err := MilestoneRepository.GetAllByRepoID(
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

	ms, err := MilestoneRepository.GetAllByRepoID(
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

	if err := MilestoneRepository.DeleteAllByRepoID(
		context.Background(),
		12,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := MilestoneRepository.GetAllByRepoID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_DeleteAllByRepoID_ErrNoDocument(t *testing.T) {
	if err := MilestoneRepository.DeleteAllByRepoID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}