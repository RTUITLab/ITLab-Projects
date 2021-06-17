package issues_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"github.com/sirupsen/logrus"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	label "github.com/ITLab-Projects/pkg/models/label"
	mm "github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/gorilla/mux"

	s "github.com/ITLab-Projects/service/api/v1/issues"
	"github.com/ITLab-Projects/service/repoimpl"
	"github.com/go-kit/kit/log"
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
		s.MakeEndPoints(service),
		Router,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_GetIssuesHTTP(t *testing.T) {
	if err := RepoImp.Issue.Save(
		context.Background(),
		[]mm.IssuesWithMilestoneID{
			{
				MilestoneID: 1,
				RepoID:      2,
				Issue: mm.Issue{
					ID:    1,
					Title: "mock_issue_1",
					State: "open",
				},
			},
			{
				MilestoneID: 1,
				RepoID:      2,
				Issue: mm.Issue{
					ID:    2,
					Title: "mock_issue_2",
					State: "open",
				},
			},
			{
				MilestoneID: 1,
				RepoID:      2,
				Issue: mm.Issue{
					ID:    3,
					Title: "mock_issue_3",
					State: "open",
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	req := httptest.NewRequest("GET", "/issues", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.Log(w.Result().StatusCode)
		t.Log(w.Body.String())
		t.FailNow()
	}

	t.Log(w.Body.String())
}

func TestFunc_GetLabelsHTTP(t *testing.T) {
	if err := RepoImp.Issue.Save(
		context.Background(),
		[]mm.IssuesWithMilestoneID {
			{
				MilestoneID: 3,
				Issue: mm.Issue{
					ID: 1,
					Labels: []label.Label{
						{
							CompactLabel: label.CompactLabel{
								Name: "mock_1",
							},
						},
						{
							CompactLabel: label.CompactLabel{
								Name: "mock_2",
							},
						},
					},
					State: "open",
				},
			},
			{
				MilestoneID: 4,
				Issue: mm.Issue{
					ID: 2,
					State: "open",
					Labels: []label.Label{
						{
							CompactLabel: label.CompactLabel{
								Name: "mock_3",
							},
						},
						{
							CompactLabel: label.CompactLabel{
								Name: "mock_4",
							},
						},
					},
				},
			},
		},	
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteAllIssuesByMilestonesID(
		context.Background(),
		[]uint64{3,4},
	)

	req := httptest.NewRequest("GET", "/issues/labels", nil)
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Log("ASsert error")
		t.FailNow()
	}

	t.Log(w.Body.String())
}