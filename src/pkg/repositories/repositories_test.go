package repositories_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ITLab-Projects/pkg/models/estimate"

	"github.com/ITLab-Projects/pkg/models/functask"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/models/realese"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repo"

	"github.com/ITLab-Projects/pkg/repositories"
)

var Repositories *repositories.Repositories
var requster githubreq.Requester

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	token, find := os.LookupEnv("ITLAB_PROJECTS_ACCESSKEY")
	if !find {
		panic("Don't find token")
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI")
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

	requster = githubreq.New(
		&githubreq.Config{
			AccessToken: token,
		},
	)

	logrus.Info(token)
}

func TestFunc_SaveRepo(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Repo.Save(repo.ToRepo(repos)); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_SaveRepoAndDeleteInfind(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(repos[0:10])

	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	if err := Repositories.Repo.SaveAndDeletedUnfind(
		ctx,
		repo.ToRepo(repos[0:10]),
	); err != nil {
		t.Log(err)
	}

	t.Log(Repositories.Repo.Count())
}

func TestFunc_GetAllRepos(t *testing.T) {
	var repos []repo.Repo
	err := Repositories.Repo.GetAll(&repos)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(repos) == 0 {
		t.Log(err)
		t.FailNow()
	}

	t.Log(len(repos))
}

func TestFunc_SaveMilestones(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var wg sync.WaitGroup
	for i, _ := range repos {
		wg.Add(1)
		go func(r *repo.RepoWithURLS, wg *sync.WaitGroup) {
			defer wg.Done()
			ms, err := requster.GetMilestonesForRepoWithID(r.Repo)
			if err != nil {
				t.Log(err)
			}
			if err := Repositories.Milestone.Save(ms); err != nil {
				t.Log(err)
			}
		}(&repos[i], &wg)
	}
	wg.Wait()
}

func TestFunc_GetAllMilestones(t *testing.T) {
	var ms []milestone.MilestoneInRepo
	err := Repositories.Milestone.GetAll(&ms)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(len(ms))
}

func TestFunc_SaveRealese(t *testing.T) {
	realse := realese.RealeseInRepo{
		RepoID: 1,
		Realese: realese.Realese{
			ID:      2,
			HTMLURL: "some_html_url",
			URL:     "some_url",
		},
	}

	err := Repositories.Realese.Save(realse)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetRealse(t *testing.T) {
	opts := options.FindOne()
	var rel realese.RealeseInRepo
	err := Repositories.Realese.GetOne(
		context.Background(),
		bson.M{"repoid": 1},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&rel)
		},
		opts,
	)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(rel)
}

func TestFunc_SaveFuncTask(t *testing.T) {
	ft := functask.FuncTask{
		MilestoneID: 2,
		FuncTaskURL: "some_url",
	}

	if err := Repositories.FuncTask.Save(ft); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetFuncTask(t *testing.T) {
	var ft functask.FuncTask

	if err := Repositories.FuncTask.GetOne(
		context.Background(),
		bson.M{"milestone_id": 2},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&ft)
		},
		options.FindOne(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(ft)
}

func TestFunc_DeleteTaskFunc(t *testing.T) {
	if err := Repositories.FuncTask.Delete(2); err != nil {
		t.Log(err)
		t.FailNow()
	}
}


func TestFunc_SaveEstimate(t *testing.T) {
	e := estimate.Estimate{
		MilestoneID: 2,
		EstimateURL: "some_url",
	}

	if err := Repositories.Estimate.Save(e); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetEstimate(t *testing.T) {
	var e estimate.Estimate

	if err := Repositories.Estimate.GetOne(
		context.Background(),
		bson.M{"milestone_id": 2},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&e)
		},
		options.FindOne(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(e)
}

func TestFunc_DeleteEstimate(t *testing.T) {
	if err := Repositories.Estimate.Delete(2); err != nil {
		t.Log(err)
		t.FailNow()
	}
}