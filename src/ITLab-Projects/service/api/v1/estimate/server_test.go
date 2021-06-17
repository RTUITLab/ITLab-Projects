package estimate_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	mm "github.com/ITLab-Projects/pkg/models/milestone"

	me "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"go.mongodb.org/mongo-driver/bson/primitive"

	kitl "github.com/go-kit/kit/log/logrus"
	"github.com/gorilla/mux"

	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/mfsreq"
	s "github.com/ITLab-Projects/service/api/v1/estimate"
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
		kitl.NewLogrusLogger(logrus.StandardLogger()),
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
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

func TestFunc_AddEstimate_ErrDontFindMilestone(t *testing.T) {
	est := me.EstimateFile{
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

	req := httptest.NewRequest("POST", "/estimate", bytes.NewBuffer(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Log("Assert Error")
		t.FailNow()
	}
	t.Log(w.Body)

}

func TestFunc_AddEstimate_(t *testing.T) {
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

	est := me.EstimateFile{
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

	req := httptest.NewRequest("POST", "/estimate", bytes.NewBuffer(data))

	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Log("Asser Error")
		t.FailNow()
	}

}

func TestFunc_DeleteEstimate_DontFindEstimate(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/estimate/1", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Log("Assert error")
		t.FailNow()
	}
	t.Log(w.Body)
}

func TestFunc_DeleteEstimate_(t *testing.T) {
	RepoImp.Estimate.Save(
		context.Background(),
		me.EstimateFile {
			milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: primitive.NewObjectID(),
			},
		},
	)
	req := httptest.NewRequest("DELETE", "/estimate/1", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)
	t.Log(w.Result().StatusCode)
	t.Log(w.Body)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.FailNow()
	}
}