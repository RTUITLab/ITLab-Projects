package landing_test

import (
	"context"
	"os"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/statuscode"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/api/v1/landing"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var service landing.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
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
	RepoImp = repoimpl.New(Repositories)

	service = landing.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("INIT")
}

func TestFunc_LandingTests(t *testing.T) {
	if err := RepoImp.LandingRepositoryImp.SaveAndDeleteUnfindLanding(
		context.Background(),
		[]model.Landing{
			{
				LandingCompact: model.LandingCompact{
					RepoId: 1,
					Title: "mock_1",
					Tags: []string{"Backend", "Web"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 2,
					Title: "mock_2",
					Tags: []string{"Web"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 3,
					Title: "mock_3",
					Tags: []string{"VR"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer Repositories.Landing.DeleteMany(
		context.Background(),
		bson.M{},
		nil,
	)
	testfunc_GetAllLandings_WithOutParams(t)
	testfunc_GetAllLandings_ByName(t)
	testfunc_GetAllLandings_ByTag(t)
	testfunc_GetByID(t)
	testfunc_GetByID_NotFound(t)
}

func testfunc_GetAllLandings_WithOutParams(t *testing.T) {
	ls, err := service.GetAllLandings(
		context.Background(),
		0,
		0,
		"",
		"",
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 3 {
		t.Log("assert err")
		t.FailNow()
	}

	for _, l := range ls {
		if !(l.Title == "mock_1" && l.RepoId == 1 || l.Title == "mock_2" && l.RepoId == 2 || l.Title == "mock_3" && l.RepoId == 3) {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func testfunc_GetAllLandings_ByName(t *testing.T) {
	ls, err := service.GetAllLandings(
		context.Background(),
		0,
		0,
		"",
		"mock_3",
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 1 {
		t.Log("Assert error")
		t.FailNow()
	}

	l := ls[0]

	if !(l.Title == "mock_3" && l.RepoId == 3) {
		t.Log("Assert error")
		t.FailNow()
	}
}

func testfunc_GetAllLandings_ByTag(t *testing.T) {
	ls, err := service.GetAllLandings(
		context.Background(),
		0,
		0,
		"Web",
		"",
	)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 2 {
		t.Log("Assert error")
		t.FailNow()
	}

	
	for _, l := range ls {
		if !(l.Title == "mock_1" && l.RepoId == 1 || l.Title == "mock_2" && l.RepoId == 2) {
			t.Log("Assert error")
			t.FailNow()
		}
	}
	
}

func testfunc_GetByID(t *testing.T) {
	l, err := service.GetLanding(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !(l.RepoId == 1 && l.Title == "mock_1") {
		t.Log("assert error")
		t.FailNow()
	}
}

func testfunc_GetByID_NotFound(t *testing.T) {
	_, err := service.GetLanding(
		context.Background(),
		123123,
	)
	if err := statuscode.GetError(err); err != landing.ErrLandingNotFound {
		t.Log(err)
		t.FailNow()
	}
}