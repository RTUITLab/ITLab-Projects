package estimate_test

import (
	"github.com/Kamva/mgm"
	"github.com/ITLab-Projects/pkg/repositories/utils/test"
	"context"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/estimate"
	"github.com/ITLab-Projects/pkg/models/milestonefile"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/estimate"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Repositories *repositories.Repositories
var EstimateRepository *estimate.EstimateRepositoryImp
func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}

	Repositories = test.GetTestRepository()

	EstimateRepository = estimate.New(
		Repositories.Estimate,
	)
	mgm.Coll(&mgm.DefaultModel{}).Database().Drop(
		context.Background(),
	)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init sucess")
}

func TestFunc_SaveEstimeate_Rewriting(t *testing.T) {
	id_1 := primitive.NewObjectID()
	id_2 := primitive.NewObjectID()

	est_1 := &model.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id_1,
		},
	}

	est_2 := &model.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id_2,
		},
	}

	if err := EstimateRepository.SaveEstimate(
		context.Background(),
		est_1,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := EstimateRepository.SaveEstimate(
		context.Background(),
		est_2,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer EstimateRepository.DeleteOneEstimateByMilestoneID(
		context.Background(),
		1,
	)

	get, err := EstimateRepository.GetEstimateByMilestoneID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if get.FileID != id_2 {
		t.Logf("Assert error: exprect %s get %s", id_2.String(), get.FileID.String())
	}
}

func TestFunc_SaveByValue(t *testing.T) {
	id := primitive.NewObjectID()
	est := model.EstimateFile{
		milestonefile.MilestoneFile{
			MilestoneID: 1,
			FileID: id,
		},
	}

	if err := EstimateRepository.SaveEstimate(
		context.Background(),
		est,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer EstimateRepository.DeleteOneEstimateByMilestoneID(
		context.Background(),
		1,
	)

	get, err := EstimateRepository.GetEstimateByMilestoneID(
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
	ests := []*model.EstimateFile{
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

	if err := EstimateRepository.SaveEstimate(
		context.Background(),
		ests,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	
	gets, err := EstimateRepository.GetEstimatesByMilestonesID(
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

	if err := EstimateRepository.DeleteManyEstimatesByMilestonesID(
		context.Background(),
		[]uint64{1,2},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if _, err := EstimateRepository.GetEstimatesByMilestonesID(
		context.Background(),
		[]uint64{1,2},
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.FailNow()
	}
}