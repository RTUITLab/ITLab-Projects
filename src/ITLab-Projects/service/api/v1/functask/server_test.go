package functask_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
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

	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/mfsreq"

	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/joho/godotenv"
)

var Router *mux.Router

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories := test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)

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
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_AddFunctask_ErrDontFindMilestone(t *testing.T) {
	est := mf.FuncTaskFile{
		MilestoneFile: milestonefile.MilestoneFile{
			FileID:      primitive.NewObjectID(),
		},
	}
	data, err := json.Marshal(est)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(string(data))

	req := httptest.NewRequest("POST", "/task/12", bytes.NewBuffer(data))

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
		},
	}
	data, err := json.Marshal(est)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(string(data))

	req := httptest.NewRequest("POST", "/task/12", bytes.NewBuffer(data))

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