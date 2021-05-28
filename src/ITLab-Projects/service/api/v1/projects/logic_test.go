package projects_test

import (
	"github.com/sirupsen/logrus"
	mre "github.com/ITLab-Projects/pkg/models/realese"
	mt "github.com/ITLab-Projects/pkg/models/tag"

	"context"
	"net/http"
	"os"
	"testing"

	me "github.com/ITLab-Projects/pkg/models/estimate"
	mf "github.com/ITLab-Projects/pkg/models/functask"
	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	mr "github.com/ITLab-Projects/pkg/models/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
		logrus.Warn("Don't find env")
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI_TEST")
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
	t.Log("Deprecated")
	t.SkipNow()
	if err := service.UpdateProjects(
		context.Background(),
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFunc_GetProjects(t *testing.T) {
	t.Log("Deprecated")
	t.SkipNow()
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
	t.Log("Deprecated")
	t.SkipNow()
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

func TestFunc_DeleteProject(t *testing.T) {
	if err := RepoImp.Repo.Save(
		context.Background(),
		mr.Repo{
			ID: 12,
			Name: "mock_repo_1",
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.Milestone.Save(
		context.Background(),
		[]mm.MilestoneInRepo{
			{	
				RepoID: 12,
				Milestone: mm.Milestone{
					MilestoneFromGH: mm.MilestoneFromGH{
						ID: 1,
						Title: "mock_milestone_1",
					},
				},
			},
			{	
				RepoID: 12,
				Milestone: mm.Milestone{
					MilestoneFromGH: mm.MilestoneFromGH{
						ID: 2,
						Title: "mock_milestone_2",
					},
				},
			},
			{	
				RepoID: 12,
				Milestone: mm.Milestone{
					MilestoneFromGH: mm.MilestoneFromGH{
						ID: 3,
						Title: "mock_milestone_3",
					},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.Issue.Save(
		context.Background(),
		[]mm.IssuesWithMilestoneID{
			{
				RepoID: 12,
				MilestoneID: 1,
				Issue: mm.Issue{
					ID: 1,
					Title: "mock_issue_1",
				},
			},
			{
				RepoID: 12,
				MilestoneID: 1,
				Issue: mm.Issue{
					ID: 2,
					Title: "mock_issue_2",
				},
			},
			{
				RepoID: 12,
				MilestoneID: 2,
				Issue: mm.Issue{
					ID: 3,
					Title: "mock_issue_3",
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.Estimate.Save(
		context.Background(),
		me.EstimateFile{
			MilestoneFile: milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: primitive.NewObjectID(),
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.FuncTask.Save(
		context.Background(),
		mf.FuncTaskFile{
			MilestoneFile: milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: primitive.NewObjectID(),
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.Tag.Save(
		context.Background(),
		[]mt.Tag{
			{
				RepoID: 12,
				Tag: "mock_tag_1",
			},
			{
				RepoID: 12,
				Tag: "mock_tag_2",
			},
			{
				RepoID: 12,
				Tag: "mock_tag_3",
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := RepoImp.Realese.Save(
		context.Background(),
		mre.RealeseInRepo {
			RepoID: 12,
			Realese: mre.Realese{
				ID: 1,
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}


	if err := service.DeleteProject(
		context.Background(),
		12,
		&http.Request{
			Header: http.Header{},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetFilteredTagsByRepoID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetAllIssuesByMilestoneID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetAllIssuesByMilestoneID(
		context.Background(),
		2,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetAllMilestonesByRepoID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetFuncTaskByMilestoneID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetEstimateByMilestoneID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetByID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}