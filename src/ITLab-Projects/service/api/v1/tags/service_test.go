package tags_test

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

	"github.com/ITLab-Projects/pkg/models/landing"
	mt "github.com/ITLab-Projects/pkg/models/tag"

	s "github.com/ITLab-Projects/service/api/v1/tags"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Router *mux.Router
func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

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
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_GetTagsHTTP(t *testing.T) {
	if err := Repositories.Landing.Save(
		context.Background(),
		[]landing.Landing{
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 1,
					Tags: []string{"mock_tag_1", "mock_tag_2"},
				},
			},
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 2,
					Tags: []string{"mock_tag_3"},
				},
			},
			{
				LandingCompact: landing.LandingCompact{
					RepoId: 4,
					Tags: []string{"mock_tag_4"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.LandingRepositoryImp.DeleteLandingsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.LandingRepositoryImp.DeleteLandingsByRepoID(
		context.Background(),
		2,
	)
	defer RepoImp.LandingRepositoryImp.DeleteLandingsByRepoID(
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
			t.Log(tag.Tag)
		default:
			t.Log("Assert error")
			t.FailNow()
		}
	}
}