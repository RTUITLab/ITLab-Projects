package githubreq_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repo"
)

// TODO rewrite test
// Need to write mock gh service

var requster *githubreq.GHRequester

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
	t.Skipf("Deprecated while refactor")
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
	t.Skipf("Deprecated while refactor")
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
	t.Skipf("Deprecated while refactor")
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

func TestFunc_GetIssues(t *testing.T) {
	t.Skipf("Deprecated while refactor")
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

	count := 0
	for _, m := range ms {
		count += len(m.Issues)
	}

	t.Log(count)
	t.Log("milestone count: ", len(ms))

	var is []milestone.IssuesWithMilestoneID

	isChan := make(chan milestone.IssuesWithMilestoneID)

	count = 0
	for j, _ := range ms {
		for i, _ := range ms[j].Issues {
			count++
			go func(i milestone.Issue, MID, RepoID uint64) {
				isChan <- milestone.IssuesWithMilestoneID{
					MilestoneID: MID,
					RepoID:      RepoID,
					Issue:       i,
				}
			}(ms[j].Issues[i], ms[j].Milestone.ID, ms[j].RepoID)
		}
	}

	for i := 0; i < count; i++ {
		select {
		case issue := <-isChan:
			is = append(is, issue)
		}
	}

	t.Log("Operations: ", count)

	t.Log(len(is))
}

func TestFunc_GetLastRealese(t *testing.T) {
	t.Skipf("Deprecated while refactor")
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	rls, err := requster.GetLastsRealeseWithRepoID(
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
	t.Log(rls)
}

func TestFunc_GetAllMilestonesForRepoWithID(t *testing.T) {
	t.Skipf("Deprecated while refactor")
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(len(repos))

	msChan := make(chan []milestone.MilestoneInRepo, 1)
	go func() {
		defer close(msChan)
		ms, err := requster.GetAllMilestonesForRepoWithID(
			context.Background(),
			repo.ToRepo(repos),
			func(e error) {
				logrus.WithFields(
					logrus.Fields{
						"err": e,
					},
				).Error()
			},
		)
		if err != nil {
			return
		}
		msChan <- ms
	}()

	t.Log(<-msChan)
}

func TestFunc_GetTagsForRepos(t *testing.T) {
	t.Skipf("Deprecated while refactor")
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

	t.Log(tags)
}

func TestFunc_HashSet(t *testing.T) {
	var issues []milestone.IssueFromGH = []milestone.IssueFromGH{
		{Issue: milestone.Issue{Description: "issue_1"}, Milestone: nil},
		{Issue: milestone.Issue{Description: "issue_2"}, Milestone: nil},
		{Issue: milestone.Issue{Description: "issue_3"}, Milestone: &milestone.MilestoneFromGH{ID: 2}},
		{Issue: milestone.Issue{Description: "issue_4"}, Milestone: &milestone.MilestoneFromGH{ID: 2}},
		{Issue: milestone.Issue{Description: "issue_5"}, Milestone: &milestone.MilestoneFromGH{ID: 1}},
		{Issue: milestone.Issue{Description: "issue_6"}, Milestone: &milestone.MilestoneFromGH{ID: 3}},
	}

	set := make(map[interface{}][]milestone.Issue)

	for _, issue := range issues {
		if issue.Milestone != nil {
			if _, find := set[*issue.Milestone]; !find {
				set[*issue.Milestone] = []milestone.Issue{issue.Issue}
			} else if find {
				set[*issue.Milestone] = append(set[*issue.Milestone], issue.Issue)
			}
		}
	}

	var milestones []milestone.Milestone

	for k, v := range set {
		m := k.(milestone.MilestoneFromGH)
		milestones = append(milestones, milestone.Milestone{MilestoneFromGH: m, Issues: v})
	}

	for _, m := range milestones {
		t.Log(m)
	}

}