package functask_test

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"os"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/functask"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/functask"
	"github.com/joho/godotenv"
)

var Repositories *repositories.Repositories
var FuncTaskRepository *functask.FuncTaskRepositoryImp
func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
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

	FuncTaskRepository = functask.New(
		_r.FuncTask,
	)
}

func TestFunc_init(t *testing.T) {
	t.Log("Init")
}

func TestFunc_SaveFuncTask_Rewriting(t *testing.T) {
	id_1 := primitive.NewObjectID()
	id_2 := primitive.NewObjectID()

	task_1 := &model.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id_1,
		},
	}

	task_2 := &model.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id_2,
		},
	}

	if err := FuncTaskRepository.SaveFuncTask(
		context.Background(),
		task_1,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := FuncTaskRepository.SaveFuncTask(
		context.Background(),
		task_2,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	defer FuncTaskRepository.DeleteOneFuncTaskByMilestoneID(
		context.Background(),
		1,
	)

	get, err := FuncTaskRepository.GetFuncTaskByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if get.FileID != id_2 {
		t.Logf("Assert error: expected %s get %s", id_2, get.FileID)
		t.FailNow()
	}
}

func TestFunc_SaveByValue(t *testing.T) {
	id := primitive.NewObjectID()
	est := model.FuncTaskFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id,
		},
	}

	if err := FuncTaskRepository.SaveFuncTask(
		context.Background(),
		est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer FuncTaskRepository.DeleteOneFuncTaskByMilestoneID(
		context.Background(),
		1,
	)

	get, err := FuncTaskRepository.GetFuncTaskByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if get.FileID != id {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_GetManyEstimates_AND_DeleteMany(t *testing.T) {
	id_1 := primitive.NewObjectID()
	id_2 := primitive.NewObjectID()
	ests := []*model.FuncTaskFile{
		{
			milestonefile.MilestoneFile{
				MilestoneID: 1,
				FileID: id_1,
			},
		},
		{
			milestonefile.MilestoneFile{
				MilestoneID: 2,
				FileID: id_2,
			},
		},
	}

	if err := FuncTaskRepository.SaveFuncTask(
		context.Background(),
		ests,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	
	gets, err := FuncTaskRepository.GetFuncTasksByMilestonesID(
		context.Background(),
		[]uint64{1,2},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, g := range gets {
		if !(g.MilestoneID == 1 && g.FileID == id_1 || g.MilestoneID == 2 && g.FileID == id_2) {
			t.Log("Assert error")
			t.FailNow()
		}
	}

	if err := FuncTaskRepository.DeleteManyFuncTasksByMilestonesID(
		context.Background(),
		[]uint64{1,2},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := FuncTaskRepository.GetFuncTasksByMilestonesID(
		context.Background(),
		[]uint64{1,2},
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}