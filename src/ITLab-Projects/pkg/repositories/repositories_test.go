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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	err := Repositories.Repo.GetAllFiltered(
		context.Background(),
		bson.M{},
		func(c *mongo.Cursor) error {
			return c.All(
				context.Background(),
				&repos,
			)
		},
		options.Find(),
	)
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

func TestFunc_SaveMilestonesAndDeleteUnfind(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	ms, err := requster.GetAllMilestonesForRepoWithID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	t.Log(ms)

	if err := Repositories.Milestone.SaveAndDeletedUnfind(
		ctx,
		ms,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(Repositories.Milestone.Count())
}

func TestFunc_GetAllMilestones(t *testing.T) {
	var ms []milestone.MilestoneInRepo
	err := Repositories.Milestone.GetAllFiltered(
		context.Background(),
		bson.M{},
		func(c *mongo.Cursor) error {
			return c.All(
				context.Background(),
				&ms,
			)
		},
		options.Find(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(len(ms))
}

func TestFunc_SaveRealese(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	realse, err := requster.GetLastsRealeseWithRepoID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Realese.Save(realse); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_SaveRealeseAndDeleteUnfind(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	realse, err := requster.GetLastsRealeseWithRepoID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)
	if err != nil {
		t.Log(err)
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := Repositories.Realese.SaveAndDeletedUnfind(
		ctx,
		realse,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetRealse(t *testing.T) {
	opts := options.FindOne()
	var rel realese.RealeseInRepo
	err := Repositories.Realese.GetOne(
		context.Background(),
		bson.M{"repoid": 174697113},
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

func TestFunc_SaveTags(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	tags, err := requster.GetAllTagsForRepoWithID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Tag.Save(tags); err != nil {
		t.Log(err)
		t.FailNow()
	}

}

func TestFunc_SaveAndDeleteUnfindTags(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	tags, err := requster.GetAllTagsForRepoWithID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := Repositories.Tag.SaveAndDeletedUnfind(
		ctx,
		tags,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_SaveAndDeleteUnfindTIssue(t *testing.T) {
	Repositories.Issue.DeleteOne(
		context.Background(),
		bson.M{"id": 1},
		nil,
		options.Delete(),
	)
	var issues []milestone.IssuesWithMilestoneID
	if err := Repositories.Issue.GetAllFiltered(
		context.Background(),
		bson.M{},
		func(c *mongo.Cursor) error {
			return c.All(
				context.Background(),
				&issues,
			)
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Issue.Save(
		milestone.IssuesWithMilestoneID{
			MilestoneID: 12,
			Issue: milestone.Issue{
				Title: "Mock-Issue",
				ID:    1,
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := Repositories.Issue.SaveAndUpdatenUnfind(
		ctx,
		issues,
		bson.M{"$set": bson.M{"deleted": true}},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	var issue milestone.IssuesWithMilestoneID

	if err := Repositories.Issue.GetOne(
		context.Background(),
		bson.M{"id": 1},
		func(sr *mongo.SingleResult) error {
			return sr.Decode(&issue)
		},
		options.FindOne(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !issue.Deleted {
		t.Log("Failed assert")
		t.FailNow()
	}
}

func TestFunc_GetMilestone(t *testing.T) {
	var m milestone.MilestoneInRepo

	if err := Repositories.Milestone.GetOne(
		context.Background(),
		bson.M{"id": 1},
		func(sr *mongo.SingleResult) error {
			if err := sr.Err(); err != nil {
				return err
			}
			return sr.Decode(&m)
		},
		options.FindOne(),
	); err != nil {
		t.Log(err)
	}

	t.Log(m)
}
