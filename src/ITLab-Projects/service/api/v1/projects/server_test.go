package projects_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	me "github.com/ITLab-Projects/pkg/models/estimate"
	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
	mt "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/ITLab-Projects/pkg/models/user"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mr "github.com/ITLab-Projects/pkg/models/repo"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/ITLab-Projects/service/repoimpl/tag"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Router *mux.Router
var Requester githubreq.Requester
func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err)
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI_TEST")
	if !find {
		panic("Don't find dburi")
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
	req := httptest.NewRequest("GET", "/", nil)
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
	req := httptest.NewRequest("GET", "/?name=mock_1", nil)
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

	if err := Repositories.Tag.Save(
		context.Background(),
		[]mt.Tag{
			{
				RepoID: 1,
				Tag: "mock_tag",
			},
			{
				RepoID: 2,
				Tag: "mock_tag",
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.TagRepositoryImp.DeleteTagsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.TagRepositoryImp.DeleteTagsByRepoID(
		context.Background(),
		2,
	)
	req := httptest.NewRequest("GET", "/?tag=mock_tag", nil)
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

	if err := Repositories.Tag.Save(
		context.Background(),
		[]mt.Tag{
			{
				RepoID: 1,
				Tag: "mock_tag_1",
			},
			{
				RepoID: 1,
				Tag: "mock_tag_2",
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

	proj := &repoasproj.RepoAsProjPointer{}

	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Log(status)
		t.Log(w.Body.String())
		t.FailNow()
	}

	if err := json.NewDecoder(w.Body).Decode(proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if proj.Repo.Name != "mock_repo" || proj.Repo.ID != 1 {
		t.Log("Assert error")
		t.FailNow()
	}

	for _, tag := range proj.Tags {
		switch tag.Tag {
		case "mock_tag_1", "mock_tag_2":
		default:
			t.Log("assert error")
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

	req = httptest.NewRequest("DELETE", "/1", nil)
	w = httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Log(status)
		t.Log(w.Body.String())
		t.FailNow()
	}
}