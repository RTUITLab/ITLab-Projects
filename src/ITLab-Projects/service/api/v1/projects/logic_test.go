package projects_test

import (
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"github.com/ITLab-Projects/pkg/models/landing"
	mre "github.com/ITLab-Projects/pkg/models/realese"
	"github.com/Kamva/mgm"
	"github.com/sirupsen/logrus"

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

	"github.com/ITLab-Projects/pkg/mfsreq"

	"github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/go-kit/kit/log"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/joho/godotenv"
)

var service projects.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

	service = projects.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init")
	mgm.Coll(&me.EstimateFile{}).Database().Drop(context.Background())
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

	if err := RepoImp.Landing.Save(
		context.Background(),
		landing.Landing{
			LandingCompact: landing.LandingCompact{
				RepoId: 12,
				Title: "mock_1",
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

	if _, err  := RepoImp.GetLandingByRepoID(
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