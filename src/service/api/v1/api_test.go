package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ITLab-Projects/pkg/githubreq"
	"github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/repoasproj"
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

func TestFunc_GetPanic(t *testing.T) {
	f := make(chan int, 1)
	s := make(chan int, 1)
	v := make(chan int, 1)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	go func() {
		time.Sleep(50 * time.Millisecond)
		f <- 2
	}()

	go func() {
		time.Sleep(60 * time.Millisecond)
		t.Log("Error")
		cancel()
		s <- 4
	}()

	go func() {
		time.Sleep(40 * time.Millisecond)
		v <- 3
	}()

	var (
		num1 *int = nil
		num2 *int = nil
		num3 *int = nil
	)

	for i := 0; i < 3; i++ {
		t.Log("Start selection")
		select {
		case <-ctx.Done():
			t.Log("Got done returning")
			return
		case _f := <-f:
			t.Log("Catch f")
			num1 = &_f
		case _s := <-s:
			t.Log("Catch s")
			num2 = &_s
		case _v := <-v:
			t.Log("Catch v")
			num3 = &_v
		}
	}

	t.Log(*num1)
	t.Log(*num2)
	t.Log(*num3)

	t.Log("Okay")
}

func TestFunc_AddFuncTask_NotFound(t *testing.T) {
	f := functask.FuncTask{
		MilestoneID: 1,
		FuncTaskURL: "some_url",
	}

	data, err := json.Marshal(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	req := httptest.NewRequest("POST", "/api/v1/projects/task", bytes.NewReader(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_AddTestFunc(t *testing.T) {
	milestone := milestone.MilestoneInRepo{
		RepoID: 12,
		Milestone: milestone.Milestone{
			MilestoneFromGH: milestone.MilestoneFromGH{
				ID: 2,
			},
		},
	}

	if err := API.Repository.Milestone.Save(milestone); err != nil {
		t.Log(err)
		t.FailNow()
	}

	f := functask.FuncTask{
		MilestoneID: 2,
		FuncTaskURL: "some_url",
	}

	data, err := json.Marshal(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	req := httptest.NewRequest("POST", "/api/v1/projects/task", bytes.NewReader(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_AddTask_BadRequest(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/projects/task", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_AddEstimate_NotFound(t *testing.T) {
	f := estimate.Estimate{
		MilestoneID: 1,
		EstimateURL: "some_url",
	}

	data, err := json.Marshal(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	req := httptest.NewRequest("POST", "/api/v1/projects/estimate", bytes.NewReader(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_AddEstimate(t *testing.T) {
	milestone := milestone.MilestoneInRepo{
		RepoID: 12,
		Milestone: milestone.Milestone{
			MilestoneFromGH: milestone.MilestoneFromGH{
				ID: 2,
			},
		},
	}

	if err := API.Repository.Milestone.Save(milestone); err != nil {
		t.Log(err)
		t.FailNow()
	}

	f := estimate.Estimate{
		MilestoneID: 2,
		EstimateURL: "some_url",
	}

	data, err := json.Marshal(f)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	req := httptest.NewRequest("POST", "/api/v1/projects/estimate", bytes.NewReader(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())

}

func TestFunc_AddEstimate_BadRequest(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/projects/estimate", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_DeleteFuncTask_NotFound(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/v1/projects/task/1", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_DeleteFuncTask(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/v1/projects/task/2", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Maybe you forget to add functask with milestone_id 2")
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_DeleteEstimate_NotFound(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/v1/projects/estimate/1", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_DeleteEstimate(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/v1/projects/estimate/2", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Maybe you forget to add estimate with milestone_id 2")
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_GetProjects(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/projects/?start=10&count=20", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)
	var projs []repoasproj.RepoAsProjCompact

	t.Log(w.Result().StatusCode)

	json.NewDecoder(w.Result().Body).Decode(&projs)

	for _, p := range projs {
		t.Log(p.Repo.CreatedAt)
	}
}

func TestFunc_GetProjectsByTag(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/projects/?start=0&count=100&tag=Mobile+Tool", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)
	var projs []repoasproj.RepoAsProjCompact

	t.Log(w.Result().StatusCode)

	json.NewDecoder(w.Result().Body).Decode(&projs)

	for _, p := range projs {
		t.Log(p)
	}

}

func TestFunc_GetProjectsByName(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/projects/?start=0&count=100&name=CyberBird", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	var projs []repoasproj.RepoAsProjCompact

	t.Log(w.Result().StatusCode)

	json.NewDecoder(w.Result().Body).Decode(&projs)

	for _, p := range projs {
		t.Log(p)
	}

	t.Log(len(projs))
}

func TestFunc_GetProject(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/projects/358569413", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

	var proj repoasproj.RepoAsProj
	if err := json.NewDecoder(w.Body).Decode(&proj); err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(proj)
}

func TestFunc_GetProject_NotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/projects/3", nil)

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}
}

func TestFunc_ParseTime(t *testing.T) {
	const l = "2019-09-27T13:46:32Z"

	parsed, err := time.Parse(time.RFC3339, l)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(parsed.String())
}
