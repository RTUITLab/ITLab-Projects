package estimate_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"net/http"
	"testing"

	kitl "github.com/go-kit/kit/log/logrus"
	"github.com/pkg/errors"

	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/sirupsen/logrus"

	me "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/mfsreq"
	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/estimate"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/joho/godotenv"
)

var service s.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
	}

	Repositories = test.GetTestRepository()
	RepoImp = repoimpl.New(Repositories)



	service = s.New(
		RepoImp,
		kitl.NewLogrusLogger(logrus.StandardLogger()),
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_AddEstimate_ErrFailedToSave_NotFoundMilestone(t *testing.T) {
	err := service.AddEstimate(
		context.Background(),
		&me.EstimateFile{
			milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID:      primitive.NewObjectID(),
			},
		},
	)
	if status, _ := statuscode.GetStatus(err); status != http.StatusNotFound {
		t.Log("Assert error")
		t.Log(status)
		t.FailNow()
	}

	if statuscode.GetError(err) != s.ErrNotFoundMilestone {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_AddEstimate(t *testing.T) {
	if err := RepoImp.Milestone.Save(
		context.Background(),
		mm.MilestoneInRepo{
			RepoID: 12,
			Milestone: mm.Milestone{
				MilestoneFromGH: mm.MilestoneFromGH{
					ID: 1,
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer func() {
		if err := RepoImp.DeleteAllMilestonesByRepoID(
			context.Background(),
			12,
		); err != nil {
			t.Log(err)
			t.FailNow()
		}
	}()

	id := primitive.NewObjectID()
	est := me.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID:      id,
		},
	}

	if err := service.AddEstimate(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	get, err := RepoImp.GetEstimateByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteOneEstimateByMilestoneID(
		context.Background(),
		1,
	)

	if get.FileID != id {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_DeleteEstimate_NotFoundEstimate(t *testing.T) {
	err := service.DeleteEstimate(
		context.Background(),
		1,
		nil,
	)

	if status, _ := statuscode.GetStatus(err); status != http.StatusNotFound {
		t.Log("Assert error")
		t.FailNow()
	}

	if statuscode.GetError(err) != s.ErrNotFoundEstimate {
		t.Log("assert error")
		t.FailNow()
	}
}

func TestFunc_DeleteEstimate(t *testing.T) {
	if err := RepoImp.Milestone.Save(
		context.Background(),
		mm.MilestoneInRepo{
			RepoID: 12,
			Milestone: mm.Milestone{
				MilestoneFromGH: mm.MilestoneFromGH{
					ID: 1,
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	)

	id := primitive.NewObjectID()
	est := me.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID:      id,
		},
	}

	if err := service.AddEstimate(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := service.DeleteEstimate(
		context.Background(),
		1,
		&http.Request{
			Header: http.Header{},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetEstimateByMilestoneID(
		context.Background(),
		1,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}

type MockMFSRequester struct {
	ErrSwithcer int
}

func (m *MockMFSRequester) NewRequests(req *http.Request) mfsreq.Requests {
	if m.ErrSwithcer == 1 {
		return &MockMFSRequests_1{}
	} else {
		return &MockMFSRequests_2{}
	}
}

func (m *MockMFSRequester) GenerateDownloadLink(ID primitive.ObjectID) string {
	return "mock_download_ling"
}

type MockMFSRequests_1 struct {

}

func (m *MockMFSRequests_1) DeleteFile(ID primitive.ObjectID) error {
	return mfsreq.NetError
}

type MockMFSRequests_2 struct {

}

func (m *MockMFSRequests_2) DeleteFile(ID primitive.ObjectID) error {
	return &mfsreq.UnexpectedCodeErr{
		Err:	errors.Wrapf(mfsreq.ErrUnexpectedCode, "%v", 12),
		Code: 	12,
	}
}

func TestFunc_DeleteEstimate_NetError(t *testing.T) {
	var _service s.Service = s.New(
		RepoImp,
		kitl.NewLogrusLogger(logrus.StandardLogger()),
		&MockMFSRequester{ErrSwithcer: 1},		
	) 
	
	if err := RepoImp.Milestone.Save(
		context.Background(),
		mm.MilestoneInRepo {
			RepoID: 12,
			Milestone: mm.Milestone{
				MilestoneFromGH: mm.MilestoneFromGH{
					ID: 1,
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	)
	est := me.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: primitive.NewObjectID(),
		},
	}

	defer RepoImp.DeleteOneEstimateByMilestoneID(
		context.Background(),
		1,
	)

	if err := service.AddEstimate(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}


	err := _service.DeleteEstimate(
		context.Background(),
		1,
		&http.Request{
			Header: http.Header{},
		},
	)

	if err == nil {
		t.Log("Error is nil")
		t.FailNow()
	}

	if status, _ := statuscode.GetStatus(err); status != http.StatusConflict {
		t.Log(status)
		t.Log("Assert error")
		t.FailNow()
	}

	getErr := statuscode.GetError(err)

	t.Log(getErr)
}

func TestFunc_DeleteEstimate_UnexcpectedCode(t *testing.T) {
	var _service s.Service = s.New(
		RepoImp,
		kitl.NewLogrusLogger(logrus.StandardLogger()),
		&MockMFSRequester{ErrSwithcer: 2},	
	) 
	
	if err := RepoImp.Milestone.Save(
		context.Background(),
		mm.MilestoneInRepo {
			RepoID: 12,
			Milestone: mm.Milestone{
				MilestoneFromGH: mm.MilestoneFromGH{
					ID: 1,
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer RepoImp.DeleteAllMilestonesByRepoID(
		context.Background(),
		12,
	)
	est := me.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: primitive.NewObjectID(),
		},
	}

	defer RepoImp.DeleteOneEstimateByMilestoneID(
		context.Background(),
		1,
	)

	if err := service.AddEstimate(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}


	err := _service.DeleteEstimate(
		context.Background(),
		1,
		&http.Request{
			Header: http.Header{},
		},
	)

	if err == nil {
		t.Log("Error is nil")
		t.FailNow()
	}

	if status, _ := statuscode.GetStatus(err); status != http.StatusConflict {
		t.Log(status)
		t.Log("Assert error")
		t.FailNow()
	}

	getErr := statuscode.GetError(err)

	t.Log(getErr)
}