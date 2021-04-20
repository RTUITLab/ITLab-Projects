package v1_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/repositories"
	v1 "github.com/ITLab-Projects/service/api/v1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var API *v1.Api
var Router *mux.Router

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
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

	requster := githubreq.New(
		&githubreq.Config{
			AccessToken: token,
		},
	)

	logrus.Info(token)

	API = &v1.Api{
		Requester:  requster,
		Repository: _r,
	}

	Router = mux.NewRouter()
	API.Build(Router)
}

func TestFunc_UpdateAllProjects(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/projects/", nil)
	
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Not okay")
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}
}
