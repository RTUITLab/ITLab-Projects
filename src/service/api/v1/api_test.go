package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

func TestFunc_GetPanic(t *testing.T) {
	f := make(chan int, 1)
	s := make(chan int, 1)
	v := make(chan int, 1)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)


	go func() {
		time.Sleep(50*time.Millisecond)
		f <- 2
	}()

	go func() {
		time.Sleep(60*time.Millisecond)
		t.Log("Error")
		cancel()
		s <- 4
	}()

	go func() {
		time.Sleep(40*time.Millisecond)
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
		case _f := <- f:
			t.Log("Catch f")
			num1 = &_f
		case _s := <- s:
			t.Log("Catch s")
			num2 = &_s
		case _v := <- v:
			t.Log("Catch v")
			num3 = &_v
		}
	}

	t.Log(*num1)
	t.Log(*num2)
	t.Log(*num3)

	t.Log("Okay")
}
