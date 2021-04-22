package githubreq_test

import (
	"context"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repo"
)

var requster githubreq.Requester

func init() {
	// Strange path but okay
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

func TestFunc_GetRepositoris(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("%v \n", repos)

	for _, r := range repos {
		t.Logf("name:%s langs: %v conts: %v\n", r.Name, r.Languages, r.Contributors)
	}
}

func TestFunc_GetMilestones(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var wg sync.WaitGroup

	for _, r := range repos {
		wg.Add(1)
		go func(name string, wg *sync.WaitGroup) {
			defer wg.Done()
			logrus.Infof("Start  name: %s", name)
			_, err := requster.GetMilestonesForRepo(name)
			if err != nil {
				t.Log(err)
				return
			}
			logrus.Infof("name in channel: %s", name)
		}(r.Name, &wg)
	}
	wg.Wait()
}

func TestFunc_GetMilestonesWithRepoID(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var wg sync.WaitGroup

	for i, _ := range repos {
		wg.Add(1)
		go func(rep *repo.RepoWithURLS, wg *sync.WaitGroup) {
			defer wg.Done()
			logrus.Infof("Start  name: %s", rep.Name)
			m, err := requster.GetMilestonesForRepoWithID(rep.Repo)
			if err != nil {
				t.Log(err)
				return
			}
			logrus.Infof("name in channel: %s", rep.Name)
			if len(m) > 0 {
				logrus.Infof("repoId: %v, milestoneRepoId: %v", rep.ID, m[0].RepoID)
			}
		}(&repos[i], &wg)
	}
	wg.Wait()
}

func TestFunc_URL(t *testing.T) {
	baseUrl := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "orgs/RTUITLab",
	}

	val := url.Values{}
	val.Add("page", "10")
	val.Add("count", "20")
	baseUrl.RawQuery = val.Encode()
	t.Log(baseUrl.String())
}


func TestFunc_GetLastRealese(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	rls := requster.GetLastsRealeseWithRepoID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)
	t.Log(rls)
}

func TestFunc_GetAllMilestonesForRepoWithID(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}


	msChan := make(chan []milestone.MilestoneInRepo, 1)
	go func() {
		defer close(msChan)
		ms := requster.GetAllMilestonesForRepoWithID(
			context.Background(),
			repo.ToRepo(repos), 
			func(e error) {
				logrus.WithFields(
					logrus.Fields{
						"err": err,
					},
				).Error()
			},
		)
		msChan <- ms
	}()

	t.Log(<-msChan)
}

func TestFunc_GetTagsForRepos(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	tags := requster.GetAllTagsForRepoWithID(
		context.Background(),
		repo.ToRepo(repos),
		func(e error) {
			t.Log(e)
		},
	)

	t.Log(tags)
}

func TestFunc(t *testing.T) {
	var i []int = nil
	var l []int

	l = append(l, i...)
}