package tags_test

import (
	"context"
	"testing"

	mt "github.com/ITLab-Projects/pkg/models/tag"


	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/ITLab-Projects/service/repoimpl/issue"
	"github.com/ITLab-Projects/service/repoimpl/milestone"
	"github.com/ITLab-Projects/service/repoimpl/reales"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/ITLab-Projects/service/repoimpl/tag"
	"github.com/go-kit/kit/log"

	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/repositories"
	s "github.com/ITLab-Projects/service/api/v1/tags"
	"github.com/ITLab-Projects/service/repoimpl"
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
		log.NewJSONLogger(os.Stdout),
	)
}

func TestFunc_GetTags(t *testing.T) {
	if err := RepoImp.Tag.Save(
		context.Background(),
		[]*mt.Tag{
			{
				RepoID: 1,
				Tag: "mock_tag_1",
			},
			{
				RepoID: 1,
				Tag: "mock_tag_2",
			},
			{
				RepoID: 2,
				Tag: "mock_tag_3",
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoImp.DeleteTagsByRepoID(
		context.Background(),
		1,
	)
	defer RepoImp.DeleteTagsByRepoID(
		context.Background(),
		2,
	)

	tgs, err := service.GetAllTags(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, tg := range tgs {
		switch tg.Tag {
		case "mock_tag_1", "mock_tag_2", "mock_tag_3":
		default:
			t.Log("Asser Error")
			t.FailNow()
		}
	}
}