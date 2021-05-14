package repo_test

import (
	"sort"
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/repo"
	"github.com/joho/godotenv"
)

var Repositories *repositories.Repositories
var RepoRepository *repo.RepoRepositoryImp
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

	RepoRepository = repo.New(_r.Repo)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init sucess")
}

func TestFunc_SaveReposAndSetDeletedUnfind(t *testing.T) {
	deleted := []*model.Repo{
		{ID: 12, Name: "mock_12"},
		{ID: 13, Name: "mock_13"},
		{ID: 14, Name: "mock_14"},
	}

	not_deleted := []*model.Repo{
		{ID: 15, Name: "mock_15"},
		{ID: 16, Name: "mock_16"},
		{ID: 17, Name: "mock_17"},
	}

	all := []*model.Repo{}
	all = append(all, deleted...)
	all = append(all, not_deleted...)

	if err := Repositories.Repo.Save(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	all_ids := []uint{12,13,14,15,16,17}
	defer Repositories.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": all_ids}},
		nil,
		options.Delete(),
	)

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		not_deleted,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	repos, err := RepoRepository.GetRepos(
		context.Background(),
		bson.M{},
		options.Find(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, r := range repos {
		if r.ID == 12 || r.ID == 13 || r.ID == 14 {
			if !r.Deleted {
				t.Log("Asserting error: not deleted")
				t.FailNow()
			}
		} else if r.Deleted {
				t.Log("Asserting error: deleted but should'nt")
				t.FailNow()
		}

	}
}

func BenchmarkSaveSlice(b *testing.B) {
	slice := []model.Repo{
		{ID: 12, Name: "mock_12"},
		{ID: 13, Name: "mock_13"},
		{ID: 14, Name: "mock_14"},
	}
	defer RepoRepository.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": []uint{12,13,14}}},
		nil,
		options.Delete(),
	)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Repositories.Repo.SaveAndDeletedUnfind(
				context.Background(),
				slice,
			)
		}
	})
}

func BenchmarkSave(b *testing.B) {
	elem := model.Repo {
		ID: 12,
		Name: "mock_12",
	}
	defer Repositories.Repo.DeleteOne(
		context.Background(),
		bson.M{"id": 12},
		nil,
		options.Delete(),
	)
	b.RunParallel(
		func(p *testing.PB) {
			for p.Next() {
				Repositories.Repo.Save(
					context.Background(),
					elem,
				)
			}
		},
	)
}

func BenchmarkSavePoint(b *testing.B) {
	elem := &model.Repo {
		ID: 12,
		Name: "mock_12",
	}
	defer Repositories.Repo.DeleteOne(
		context.Background(),
		bson.M{"id": 12},
		nil,
		options.Delete(),
	)
	b.RunParallel(
		func(p *testing.PB) {
			for p.Next() {
				Repositories.Repo.Save(
					context.Background(),
					elem,
				)
			}
		},
	)
}

func BenchmarkSaveSlicePointer(b *testing.B) {
	slice := []*model.Repo{
		{ID: 12, Name: "mock_12"},
		{ID: 13, Name: "mock_13"},
		{ID: 14, Name: "mock_14"},
	}
	defer RepoRepository.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": []uint{12,13,14}}},
		nil,
		options.Delete(),
	)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Repositories.Repo.SaveAndDeletedUnfind(
				context.Background(),
				slice,
			)
		}
	})
}

func TestFunc_GetFilteredRepos(t *testing.T) {
	all := []*model.Repo {
		{ID: 12, Name: "mock"},
		{ID: 13, Name: "mock"},
		{ID: 14, Name: "mock_1"},
	}

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoRepository.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": []uint{12,13,14}}},
		nil,
		options.Delete(),
	)

	repos, err := RepoRepository.GetFilteredRepos(
		context.Background(),
		bson.M{"name": "mock"},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, r := range repos {
		if r.Name != "mock" {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func TestFunc_DeleteByID(t *testing.T) {
	all := []*model.Repo {
		{ID: 12, Name: "mock"},
		{ID: 13, Name: "mock"},
		{ID: 14, Name: "mock_1"},
	}

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoRepository.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": []uint{13,14}}},
		nil,
		options.Delete(),
	)

	if err := RepoRepository.DeleteByID(
		context.Background(),
		12,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}

	repos, err := RepoRepository.GetRepos(
		context.Background(),
		bson.M{},
		options.Find(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, r := range repos {
		if r.ID == 12 {
			t.Log("Assert error: not deleted")
			t.FailNow()
		}
	}
}

func TestFunc_DeleteByID_NoDocumentsInresult(t *testing.T) {
	if err := RepoRepository.DeleteByID(
		context.Background(),
		12,
	); err != mongo.ErrNoDocuments {
		t.Log(err)
		t.Log("Assert result")
		t.FailNow()
	}
}

func TestFunc_GetByID(t *testing.T) {
	all := []*model.Repo {
		{ID: 12, Name: "mock_12"},
		{ID: 13, Name: "mock_13"},
		{ID: 14, Name: "mock_14"},
	}

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer RepoRepository.Repo.DeleteMany(
		context.Background(),
		bson.M{"id": bson.M{"$in": []uint{12, 13, 14}}},
		nil,
		options.Delete(),
	)

	repo, err := RepoRepository.GetByID(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if repo.Name != "mock_12" || repo.ID != 12 {
		t.Log("Asserting error")
		t.FailNow()
	}
	
}

func TestFunc_GetFiltrSortRepos(t *testing.T) {
	all := []*model.Repo{
		{ID: 123, Name: "mock_123"},
		{ID: 150, Name: "mock_150"},
		{ID: 3, Name: "mock_3"},
	}

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
	}

	defer Repositories.Repo.DeleteMany(
		context.Background(),
		bson.M{},
		nil,
		options.Delete(),
	)

	repos, err := RepoRepository.GetFiltrSortRepos(
		context.Background(),
		bson.M{},
		bson.D{ {"id", -1} },
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if sorted := sort.SliceIsSorted(
		repos,
		func(i, j int) bool {
			return repos[i].ID > repos[j].ID
		},
	); !sorted {
		t.Log("Assert error: not sorted")
		t.FailNow()
	}
}

func TestFunc_GetFiltrSortFromToRepos(t *testing.T) {
	all := []*model.Repo{
		{ID: 123, Name: "mock_123"},
		{ID: 150, Name: "mock_150"},
		{ID: 3, Name: "mock_3"},
		{ID: 1, Name: "mock_1"},
	}

	if err := RepoRepository.SaveReposAndSetDeletedUnfind(
		context.Background(),
		all,
	); err != nil {
		t.Log(err)
	}

	defer Repositories.Repo.DeleteMany(
		context.Background(),
		bson.M{},
		nil,
		options.Delete(),
	)

	repos, err := RepoRepository.GetFiltrSortFromToRepos(
		context.Background(),
		bson.M{},
		bson.D{ {"id", -1} },
		1,
		2,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if sorted := sort.SliceIsSorted(
		repos,
		func(i, j int) bool {
			return repos[i].ID > repos[j].ID
		},
	); !sorted {
		t.Log("Assert error: not sorted")
		t.FailNow()
	}

	if len(repos) != 2 {
		t.Log("assert error; len is not 2")
		t.FailNow()
	}

	if repos[0].ID != 123 || repos[1].ID != 3 {
		t.Log("Assert error")
		t.FailNow()
	}
}