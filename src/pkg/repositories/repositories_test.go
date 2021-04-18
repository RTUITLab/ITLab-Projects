package repositories_test

import (
	"context"
	"os"
	"sync"
	"testing"

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
var requster *githubreq.GHRequester

func init() {
	_r, err := repositories.New(&repositories.Config{
		DBURI: "mongodb://root:root@127.0.0.1:27100/ITLabProjects",
	})
	if err != nil {
		panic(err)
	}

	Repositories = _r

	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	token, find := os.LookupEnv("ITLAB_PROJECTS_ACCESSKEY")
	if !find {
		panic("Don't find token")
	}

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

	if err := Repositories.ReposRepositorier.Save(repo.ToRepo(repos)); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetAllRepos(t *testing.T) {
	var repos []repo.Repo
	err := Repositories.ReposRepositorier.GetAll(&repos)
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
			if err := Repositories.Milestoner.Save(ms); err != nil {
				t.Log(err)
			}
		}(&repos[i], &wg)
	}
	wg.Wait()
}

func TestFunc_GetAllMilestones(t *testing.T) {
	var ms []milestone.MilestoneInRepo
	err := Repositories.Milestoner.GetAll(&ms)
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

	err := Repositories.Realeser.Save(realse)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetRealse(t *testing.T) {
	opts := options.FindOne()
	var rel realese.RealeseInRepo
	err := Repositories.Realeser.GetOne(
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
