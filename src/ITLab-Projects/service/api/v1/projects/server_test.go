package projects_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	me "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/landing"
	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	"github.com/ITLab-Projects/pkg/models/user"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mr "github.com/ITLab-Projects/pkg/models/repo"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	s "github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Router *mux.Router
var Requester githubreq.Requester
func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	token, find := os.LookupEnv("ITLAB_PROJECTS_ACCESSKEY")
	if !find {
		panic("Don't find token in env")
	}

	Requester = githubreq.New(
		&githubreq.Config{
			AccessToken: token,
		},
	)

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

	service = s.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
		Requester,
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
		nil,
	)

	Router = mux.NewRouter()

	s.NewHTTPServer(
		context.Background(),
		s.MakeEndpoints(service),
		Router,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_GetProjects_HTTP(t *testing.T) {
	if err := Repositories.Repo.Save(
		context.Background(),
		[]mr.Repo{
			{
				ID: 1,
				Name: "mock_1",
			},
			{
				ID: 2,
				Name: "mock_2",
			},
			{
				ID: 3,
				Name: "mock_3",
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteByID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		2,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		3,
	)
	req := httptest.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log(w.Result().StatusCode)
		t.Log(w.Body.String())
		t.FailNow()
	}

	// t.Log(w.Body.String())

	proj := []*repoasproj.RepoAsProjCompactPointers{}
	if err := json.NewDecoder(w.Result().Body).Decode(&proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(proj) != 3 {
		t.Log("Assert error ")
		t.Log(len(proj))
		t.FailNow()
	}

	for _, p := range proj {
		switch p.Repo.Name {
		case "mock_1", "mock_2", "mock_3":
		default:
			t.Log(p.Repo.Name)
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetProjects_HTTP_ByName(t *testing.T) {
	if err := Repositories.Repo.Save(
		context.Background(),
		[]mr.Repo{
			{
				ID: 1,
				Name: "mock_1",
			},
			{
				ID: 2,
				Name: "mock_2",
			},
			{
				ID: 3,
				Name: "mock_3",
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteByID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		2,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		3,
	)
	req := httptest.NewRequest("GET", "/projects?name=mock_1", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log(w.Result().StatusCode)
		t.Log(w.Body.String())
		t.FailNow()
	}

	// t.Log(w.Body.String())

	proj := []*repoasproj.RepoAsProjCompactPointers{}
	if err := json.NewDecoder(w.Result().Body).Decode(&proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(proj) != 1 {
		t.Log("Assert error ")
		t.Log(len(proj))
		t.FailNow()
	}

	for _, p := range proj {
		switch p.Repo.Name {
		case "mock_1":
		default:
			t.Log(p.Repo.Name)
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetProjects_HTTP_ByTag(t *testing.T) {
	if err := Repositories.Repo.Save(
		context.Background(),
		[]mr.Repo{
			{
				ID: 1,
				Name: "mock_1",
			},
			{
				ID: 2,
				Name: "mock_2",
			},
			{
				ID: 3,
				Name: "mock_3",
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteByID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		2,
	)
	defer RepoImp.DeleteByID(
		context.Background(),
		3,
	)


	// {
			// 	RepoID: 1,
			// 	Tag: "mock_tag",
			// },
			// {
			// 	RepoID: 2,
			// 	Tag: "mock_tag",
			// },

	if err := Repositories.Landing.Save(
		context.Background(),
		[]landing.Landing{
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 1,
					Tags: []string{"mock_tag"},
				},
			},
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 2,
					Tags: []string{"mock_tag"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteLandingsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteLandingsByRepoID(
		context.Background(),
		2,
	)
	req := httptest.NewRequest("GET", "/projects?tag=mock_tag", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log(w.Result().StatusCode)
		t.Log(w.Body.String())
		t.FailNow()
	}

	// t.Log(w.Body.String())

	proj := []*repoasproj.RepoAsProjCompactPointers{}
	if err := json.NewDecoder(w.Result().Body).Decode(&proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(proj) != 2 {
		t.Log("Assert error ")
		t.Log(len(proj))
		t.FailNow()
	}

	for _, p := range proj {
		switch p.Repo.Name {
		case "mock_1", "mock_2":
		default:
			t.Log(p.Repo.Name)
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_GetProject_HTTP_AndDeleteThem(t *testing.T) {
	if err := Repositories.Repo.Save(
		context.Background(),
		mr.Repo{
			ID: 1,
			Name: "mock_repo",
			Contributors: []user.User{
				{
					ID: 2,
					Login: "mock_user_1",
				},
			},

		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteByID(
		context.Background(),
		1,
	)

	if err := Repositories.Milestone.Save(
		context.Background(),
		[]mm.MilestoneInRepo{
			{
				RepoID: 1,
				Milestone: mm.Milestone{
					MilestoneFromGH: mm.MilestoneFromGH{
						ID: 1,
						Title: "mock_milestone_1",
					},
				},
			},
			{
				RepoID: 1,
				Milestone: mm.Milestone{
					MilestoneFromGH: mm.MilestoneFromGH{
						ID: 2,
						Title: "mock_milestone_2",
					},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Issue.Save(
		context.Background(),
		[]mm.IssuesWithMilestoneID{
			{
				MilestoneID: 1,
				RepoID: 1,
				Issue: mm.Issue{
					ID: 1,
					Title: "mock_issue_1",
					State: "open",
				},
			},
			{
				MilestoneID: 1,
				RepoID: 1,
				Issue: mm.Issue{
					ID: 2,
					Title: "mock_issue_2",
					State: "close",
				},
			},
			{
				MilestoneID: 2,
				RepoID: 1,
				Issue: mm.Issue{
					ID: 3,
					Title: "mock_issue_3",
					State: "open",
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	
	estimate_id := primitive.NewObjectID()

	if err := Repositories.Estimate.Save(
		context.Background(),
		me.EstimateFile{
			MilestoneFile: milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: estimate_id,
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Landing.Save(
		context.Background(),
		landing.Landing{
			LandingCompact: landing.LandingCompact{
				RepoId: 1,
				Tags: []string{"mock_tag_1", "mock_tag_2"},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	proj := &repoasproj.RepoAsProjPointer{}

	req := httptest.NewRequest("GET", "/projects/1", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Log(status)
		t.Log(w.Body.String())
		t.FailNow()
	}

	t.Log(w.Body.String())
	
	if err := json.NewDecoder(w.Body).Decode(proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if proj.Repo.Name != "mock_repo" || proj.Repo.ID != 1 {
		t.Log("Assert error")
		t.FailNow()
	}

	for _, tag := range proj.Tags {
		if !(tag.Tag == "mock_tag_1" || tag.Tag == "mock_tag_2") {
			t.Log("Assert error")
			t.FailNow()
		}
	}

	for _, m := range proj.Milestones {
		switch m.Title {
		case "mock_milestone_1", "mock_milestone_2":
		default:
			t.Log("Assert error")
			t.FailNow()
		}

		for _, i := range m.Issues {
			if 	!(	i.Title == "mock_issue_1" && i.State == "open" 	||
					i.Title == "mock_issue_2" && i.State == "close" ||
					i.Title == "mock_issue_3" && i.State == "open"		) {
						t.Log("Assert error")
						t.FailNow()
				}
		}
	}

	req = httptest.NewRequest("DELETE", "/projects/1", nil)
	w = httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Log(status)
		t.Log(w.Body.String())
		t.FailNow()
	}

}