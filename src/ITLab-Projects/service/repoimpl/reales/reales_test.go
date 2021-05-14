package reales_test

import (
	"context"
	"os"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/realese"
	"github.com/stretchr/testify/assert"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/joho/godotenv"
)

var Repositories *repositories.Repositories
var RealeseRepository *reales.RealeseRepositoryImp
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

	RealeseRepository = reales.New(_r.Realese)
}

func TestFunc_Init(t *testing.T) {
	t.Log("INIT")
}

func TestFunc_SaveReales(t *testing.T) {
	realeses := []*model.RealeseInRepo{
		{RepoID: 12, Realese: model.Realese{URL: "mock_12"}},
		{RepoID: 13, Realese: model.Realese{URL: "mock_13"}},
	}

	if err := RealeseRepository.SaveRealeses(
		context.Background(),
		realeses,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	rls_1, err := RealeseRepository.GetByRepoID(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if eq := assert.ObjectsAreEqualValues(realeses[0], rls_1); !eq {
		t.Log("Assert err")
		t.FailNow()
	}
}