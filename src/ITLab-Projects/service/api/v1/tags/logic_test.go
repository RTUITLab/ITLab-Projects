package tags_test

import (
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"testing"

	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/Kamva/mgm"

	"github.com/go-kit/kit/log"

	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/tags"
	"github.com/ITLab-Projects/service/repoimpl"
)

var service s.Service

var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

	service = s.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
	)
	
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_GetTags(t *testing.T) {
	if err := RepoImp.Landing.Save(
		context.Background(),
		[]landing.Landing{
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 1,
					Tags: []string{"mock_tag_1", "mock_tag_2"},
				},
			},
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 2,
					Tags: []string{"mock_tag_3"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteLandingsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteLandingsByRepoID(
		context.Background(),
		2,
	)

	tgs, err := service.GetAllTags(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, tg := range tgs {
		switch tg.Tag {
		case "mock_tag_1", "mock_tag_2", "mock_tag_3":
		default:
			t.Log("Asser Error")
			t.Log(tg.Tag)
			t.FailNow()
		}
	}
}