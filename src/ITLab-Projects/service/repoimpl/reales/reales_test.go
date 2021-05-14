package reales_test

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/ITLab-Projects/pkg/models/realese"


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

func TestFunc_SaveReales_AndGetByRepoID(t *testing.T) {
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

	defer Repositories.Realese.DeleteMany(
		context.Background(),
		bson.M{},
		func(dr *mongo.DeleteResult) error {
			return nil
		},
		options.Delete(),
	)

	rls_1, err := RealeseRepository.GetByRepoID(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	rls_2, err := RealeseRepository.GetByRepoID(
		context.Background(),
		13,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	

	if rls_1.RepoID != 12 || rls_1.URL != "mock_12" {
		t.Log("assert error")
		t.FailNow()
	}

	if rls_2.RepoID != 13 || rls_2.URL != "mock_13" {
		t.Log("assert error")
		t.FailNow()
	}
}