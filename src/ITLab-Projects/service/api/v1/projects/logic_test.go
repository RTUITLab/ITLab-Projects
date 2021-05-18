package projects_test

import (
	"context"
	"os"
	"testing"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"

	"github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/go-kit/kit/log"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/ITLab-Projects/service/repoimpl/tag"
	"github.com/joho/godotenv"
)

var service projects.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err)
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI")
	if !find {
		panic("Don't find dburi")
	}

	token, find := os.LookupEnv("ITLAB_PROJECTS_ACCESSKEY")
	if !find {
		panic("Don't find token")
	}

	requster := githubreq.New(
		&githubreq.Config{
			AccessToken: token,
		},
	)

	_r, err := repositories.New(&repositories.Config{
		DBURI: dburi,
	})
	if err != nil {
		panic(err)
	}

	Repositories = _r
	RepoImp = &repoimpl.RepoImp{
		estimate.New(Repositories.Estimate),
		issue.New(Repositories.Issue),
		functask.New(Repositories.FuncTask),
		milestone.New(Repositories.Milestone),
		reales.New(Repositories.Realese),
		repo.New(Repositories.Repo),
		tag.New(Repositories.Tag),
	}

	service = projects.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
		requster,
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
		nil,
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init")
}

func TestFunc_UpdateAllProjects(t *testing.T) {
	if err := service.UpdateProjects(
		context.Background(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetProjects(t *testing.T) {
	projs, err := service.GetProjects(
		context.Background(),
		0,
		1000,
		"",
		"Mobile",
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, p := range projs {
		t.Logf("%v", p.Repo.Name)
		t.Logf("%v", p.Tags)
	}

	t.Log(len(projs))
}

func TestFunc_GetProject(t *testing.T) {
	proj, err := service.GetProject(
		context.Background(),
		356562826,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(proj.Completed)
	for _, m := range proj.Milestones {
		t.Log(len(m.Issues))
	}
}