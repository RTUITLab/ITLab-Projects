package functask_test

import (
	"github.com/pkg/errors"

	"context"
	"net/http"
	"os"
	"testing"

	kitl "github.com/go-kit/kit/log/logrus"

	"github.com/ITLab-Projects/pkg/statuscode"

	"github.com/ITLab-Projects/pkg/mfsreq"
	mf "github.com/ITLab-Projects/pkg/models/functask"
	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/functask"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/ITLab-Projects/service/repoimpl/tag"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var service s.Service
var Repositories *repositories.Repositories
var RepoImp *repoimpl.RepoImp

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
		mfsreq.New(
			&mfsreq.Config{
				BaseURL:  "mfs_url",
				TestMode: true,
			},
		),
		log.NewJSONLogger(os.Stdout),
	)
}

func TestFunc_AddFuncTask_ErrFailedToSave_NotFoundMilestone(t *testing.T) {
	err := service.AddFuncTask(
		context.Background(),
		&mf.FuncTaskFile{
			milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID:      primitive.NewObjectID(),
			},
		},
	)
	if status, _ := statuscode.GetStatus(err); status != http.StatusNotFound {
		t.Log("Assert error")
		t.FailNow()
	}

	if statuscode.GetError(err) != s.ErrNotFoundMilestone {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_AddFuncTask(t *testing.T) {
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
	ft := mf.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID:      id,
		},
	}

	if err := service.AddFuncTask(
		context.Background(),
		&ft,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	get, err := RepoImp.GetFuncTaskByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteOneFuncTaskByMilestoneID(
		context.Background(),
		1,
	)

	if get.FileID != id {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_DeleteFuncTask_NotFoundEstimate(t *testing.T) {
	err := service.DeleteFuncTask(
		context.Background(),
		1,
		nil,
	)

	if status, _ := statuscode.GetStatus(err); status != http.StatusNotFound {
		t.Log("Assert error")
		t.FailNow()
	}

	if statuscode.GetError(err) != s.ErrNotFound {
		t.Log("assert error")
		t.FailNow()
	}
}

func TestFunc_DeleteFuncTask(t *testing.T) {
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

	id := primitive.NewObjectID()
	est := mf.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id,
		},
	}

	if err := service.AddFuncTask(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := service.DeleteFuncTask(
		context.Background(),
		1,
		&http.Request{
			Header: http.Header{},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := RepoImp.GetFuncTaskByMilestoneID(
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


func TestFunc_DeleteFuncTask_NetError(t *testing.T) {
	var _service s.Service = s.New(
		RepoImp,
		&MockMFSRequester{ErrSwithcer: 1},
		kitl.NewLogrusLogger(logrus.StandardLogger()),		
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
	est := mf.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: primitive.NewObjectID(),
		},
	}

	defer RepoImp.DeleteOneFuncTaskByMilestoneID(
		context.Background(),
		1,
	)

	if err := service.AddFuncTask(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}


	err := _service.DeleteFuncTask(
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

func TestFunc_DeleteFuncTask_UnexcpectedCode(t *testing.T) {
	var _service s.Service = s.New(
		RepoImp,
		&MockMFSRequester{ErrSwithcer: 2},
		kitl.NewLogrusLogger(logrus.StandardLogger()),		
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
	est := mf.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: primitive.NewObjectID(),
		},
	}

	defer RepoImp.DeleteOneFuncTaskByMilestoneID(
		context.Background(),
		1,
	)

	if err := service.AddFuncTask(
		context.Background(),
		&est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}


	err := _service.DeleteFuncTask(
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