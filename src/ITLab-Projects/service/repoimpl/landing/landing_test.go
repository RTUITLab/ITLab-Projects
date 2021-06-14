package landing_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/ITLab-Projects/pkg/models/landing"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/landing"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var Repositories *repositories.Repositories
var LandingRepository *landing.LandingRepositoryImp
func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Info("Don't find env")
	}

	Repositories = test.GetTestRepository()

	LandingRepository = landing.New(
		Repositories.Landing,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("INIT")
}

func TestFunc_SaveAndDeleteUnfindLanding(t *testing.T) {
	ls_all := []*model.Landing{
		{
			LandingCompact: model.LandingCompact{
				RepoId: 1,
				Title: "mock_1",
				Date: model.Time{
					Time: time.Now(),
				},
			},
		},
		{
			LandingCompact: model.LandingCompact{
				RepoId: 2,
				Title: "mock_2",
			},
		},
		{
			LandingCompact: model.LandingCompact{
				RepoId: 3,
				Title: "mock_3",
			},
		},
	}

	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		ls_all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer Repositories.Landing.DeleteMany(
		context.Background(),
		bson.M{"repo_id": bson.M{"$in": []uint64{1,2,3}}},
		nil,
		options.Delete(),
	)

	var get_ls []*model.Landing

	if err := Repositories.Landing.GetAllFiltered(
		context.Background(),
		bson.M{},
		func(c *mongo.Cursor) error {
			return c.All(
				context.Background(),
				&get_ls,
			)
		},
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, l := range get_ls {
		switch l.Title {
		case "mock_1", "mock_2", "mock_3":
			t.Log(l.Date)
		default:
			t.Log("assert error")
			t.FailNow()
		}
	}

	ls_new := []*model.Landing{
		{
			LandingCompact: model.LandingCompact{
				RepoId: 2,
				Title: "mock_2",
			},
		},
		{
			LandingCompact: model.LandingCompact{
				RepoId: 3,
				Title: "mock_3",
			},
		},
	}

	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		ls_new,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	var get []*model.Landing

	if err := Repositories.Landing.GetAllFiltered(
		context.Background(),
		bson.M{},
		func(c *mongo.Cursor) error {
			return c.All(
				context.Background(),
				&get,
			)
		},
		options.Find(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, l := range get {
		switch l.Title {
		case "mock_2", "mock_3":
			t.Log(l.Date)
		case "mock_1":
			t.Log("Assert error")
		default:
			t.Log("assert error")
			t.FailNow()
		}
	}
}

func TestFunc_DeleteLandingByRepoID(t *testing.T) {
	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		model.Landing{
			LandingCompact: model.LandingCompact{
				RepoId: 1,
				Title: "mock_1",
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	get, err := LandingRepository.GetLandingByRepoID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !(get.RepoId == 1 && get.Title == "mock_1") {
		t.Log("assert error")
		t.FailNow()
	}

	if err := LandingRepository.DeleteLandingsByRepoID(
		context.Background(),
		1,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := LandingRepository.GetLandingByRepoID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log("Assert err")
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetLandingTagsByRepoID(t *testing.T) {
	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		&model.Landing{
			LandingCompact: model.LandingCompact{
				RepoId: 1,
				Title: "mock_1",
				Tags: []string{"Backend", "VR", "AR"},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	tags, err := LandingRepository.GetLandingTagsByRepoID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, tag := range tags {
		switch tag.Tag {
		case "Backend", "VR", "AR":
			t.Log(tag.Tag)
		default:
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetIdsOfReposByTags(t *testing.T) {
	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		[]model.Landing{
			{
				LandingCompact: model.LandingCompact{
					RepoId: 1,
					Title: "mock_1",
					Tags: []string{"Backend", "AR", "VR"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 2,
					Title: "mock_2",
					Tags: []string{"AR", "VR"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 3,
					Title: "mock_1",
					Tags: []string{"Go"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	ids, err := LandingRepository.GetIDsOfReposByLandingTags(
		context.Background(),
		[]string{"Backend", "AR", "Go", "VR"},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range ids {
		switch i {
		case 1, 2, 3:
		default:
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetAllTags(t *testing.T) {
	if err := LandingRepository.SaveAndDeleteUnfindLanding(
		context.Background(),
		[]model.Landing{
			{
				LandingCompact: model.LandingCompact{
					RepoId: 1,
					Tags: []string{"VR", "AR"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 2,
					Tags: []string{"Backend", "AR"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 3,
					Tags: []string{"VR", "Web"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	tags, err := LandingRepository.GetAllTags(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(tags) == 0 {
		t.Log("Assert error")
		t.FailNow()
	}

	if len(tags) != 4 {
		t.Log("assert error")
		t.FailNow()
	}

	// TODO доделать тест
}