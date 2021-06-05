package landing_test

import (
	"context"
	"os"
	"testing"
	"time"

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

	LandingRepository = landing.New(
		Repositories.Landing,
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


}