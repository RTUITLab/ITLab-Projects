package tags_test

import (
	"github.com/sirupsen/logrus"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mt "github.com/ITLab-Projects/pkg/models/tag"

	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/tags"
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
func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI_TEST")
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
	)

	Router = mux.NewRouter()

	s.NewHTTPServer(
		context.Background(),
		s.MakeEndpoints(service),
		Router,
	)
}

func TestFunc_GetTagsHTTP(t *testing.T) {
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
			{
				RepoID: 2,
				Tag: "mock_tag_3",
			},
			{
				RepoID: 4,
				Tag: "mock_tag_4",
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.DeleteTagsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteTagsByRepoID(
		context.Background(),
		2,
	)
	defer RepoImp.DeleteTagsByRepoID(
		context.Background(),
		4,
	)

	req := httptest.NewRequest("GET", "/tags", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Log(status)
		t.Log(w.Body.String())
		t.FailNow()
	}

	var tags []*mt.Tag

	if err := json.NewDecoder(w.Body).Decode(&tags); err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, tag := range tags {
		switch tag.Tag {
		case "mock_tag_1", "mock_tag_2", "mock_tag_3", "mock_tag_4":
		default:
			t.Log("Assert error")
			t.FailNow()
		}
	}
}