package issues_test

import (
	"github.com/sirupsen/logrus"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ITLab-Projects/pkg/models/label"
	mm "github.com/ITLab-Projects/pkg/models/milestone"

	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/issues"
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
)

var service s.Service
var Repositories *repositories.Repositories
var RepoImp	*repoimpl.RepoImp

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		logrus.Warn("Don't find env")
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
		log.NewJSONLogger(os.Stdout),
	)
}

func TestFunc_GetIssues(t *testing.T) {
	is, err := service.GetIssues(
		context.Background(),
		0,
		10000,
		"Orbital 360 Model",
		"",
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, i := range is {
		t.Log(i.CreatedAt)
	}

	t.Log(len(is))
}

func TestFunc_GetLabels(t *testing.T) {
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
				},
			},
			{
				MilestoneID: 4,
				Issue: mm.Issue{
					ID: 2,
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

	labels, err := service.GetLabels(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, l := range labels {
		switch fmt.Sprint(l) {
		case "mock_1", "mock_2", "mock_3", "mock_4":
		default:
			t.Log(fmt.Sprint(l))
			t.Log("Assert error")
			t.FailNow()
		}
	}
}