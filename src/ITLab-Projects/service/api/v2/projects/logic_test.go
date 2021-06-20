package projects_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"github.com/Kamva/mgm"
	"github.com/gorilla/mux"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/api/v2/projects"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var service projects.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp
var Router	*mux.Router

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

	service = projects.New(
		RepoImp,
		log.NewJSONLogger(os.Stdout),
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)

	Router = mux.NewRouter()
	projects.NewHTTPServer(
		context.Background(),
		projects.MakeEndpoints(service),
		Router,
	)
}

func TestFunc_GetProjects_HTTP(t *testing.T) {
	for i := 0; i < 100; i++ {
		Repositories.Repo.Save(
			context.Background(),
			repo.Repo{
				ID: uint64(i),
				Name: fmt.Sprintf("mock_%v", i),
			},
		)
	}

	req := httptest.NewRequest(
		"GET",
		"/projects?start=20&count=60",
		nil,
	)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(
		resp,
		req,
	)
}

