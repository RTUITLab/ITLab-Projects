package functask_test

import (
	s "github.com/ITLab-Projects/service/api/v1/functask"
	mf "github.com/ITLab-Projects/pkg/models/functask"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	mm "github.com/ITLab-Projects/pkg/models/milestone"

	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"go.mongodb.org/mongo-driver/bson/primitive"

	kitl "github.com/go-kit/kit/log/logrus"
	"github.com/gorilla/mux"

	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/mfsreq"
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

var Router *mux.Router

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err)
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

	Repositories := _r
	RepoImp := &repoimpl.RepoImp{
		estimate.New(Repositories.Estimate),
		issue.New(Repositories.Issue),
		functask.New(Repositories.FuncTask),
		milestone.New(Repositories.Milestone),
		reales.New(Repositories.Realese),
		repo.New(Repositories.Repo),
		tag.New(Repositories.Tag),
	}

	service := s.New(
		RepoImp,
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
		kitl.NewLogrusLogger(logrus.StandardLogger()),
	)

	Router = mux.NewRouter()

	s.NewHTTPServer(
		context.Background(),
		s.MakeEndPoints(service),
		Router,
	)
}

func TestFunc_AddFunctask_ErrDontFindMilestone(t *testing.T) {
	est := mf.FuncTaskFile{
		MilestoneFile: milestonefile.MilestoneFile{
			FileID:      primitive.NewObjectID(),
			MilestoneID: 12,
		},
	}
	data, err := json.Marshal(est)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(string(data))

	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log("Assert Error")
		t.FailNow()
	}
	t.Log(w.Body)
}

func TestFunc_AddFuncTask_(t *testing.T) {
	if err := RepoImp.Milestone.Save(
		context.Background(),
		mm.MilestoneInRepo{
			RepoID: 1,
			Milestone: mm.Milestone{
				MilestoneFromGH: mm.MilestoneFromGH{
					ID: 12,
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.DeleteAllMilestonesByRepoID(
		context.Background(),
		1,
	)

	est := mf.FuncTaskFile{
		MilestoneFile: milestonefile.MilestoneFile{
			FileID:      primitive.NewObjectID(),
			MilestoneID: 12,
		},
	}
	data, err := json.Marshal(est)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(string(data))

	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Log("Asser Error")
		t.Log(w.Result().StatusCode)
		t.FailNow()
	}

}

func TestFunc_DeleteFuncTask_DontFindFuncTask(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/task/1", nil)
	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Log("Assert error")
		t.FailNow()
	}
	t.Log(w.Body)
}

func TestFunc_DeleteFuncTask_(t *testing.T) {
	RepoImp.FuncTask.Save(
		context.Background(),
		mf.FuncTaskFile{
			milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: primitive.NewObjectID(),
			},
		},
	)

	req := httptest.NewRequest("DELETE", "/task/1", nil)
	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.FailNow()
	}
}