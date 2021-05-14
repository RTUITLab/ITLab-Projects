package tag_test

import (
	"context"
	"os"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/tag"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/service/repoimpl/tag"
	"github.com/joho/godotenv"
)

var Repositories *repositories.Repositories
var TagRepository *tag.TagRepositoryImp
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

	TagRepository = tag.New(Repositories.Tag)
}

func TestFunc_Init(t *testing.T) {
	t.Log("Init sucess")
}

func TestFunc_GetAllTags(t *testing.T) {
	if err := Repositories.Tag.Save(
		context.Background(),
		&model.Tag{
			RepoID: 12,
			Tag: "MOCK_TAG",
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer TagRepository.DeleteTagsByRepoID(
		context.Background(),
		12,
	)

	tags, err := TagRepository.GetAllTags(
		context.Background(),
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Contains(
		t,
		tags,
		&model.Tag {
			Tag: "MOCK_TAG",
		},
	)
}

func TestFunc_GetFilteredByRepoID(t *testing.T) {
	excpected := []*model.Tag{
		{RepoID: 12, Tag: "Game"},
		{RepoID: 12, Tag: "Go"},
	}

	_tags := []*model.Tag{
		{RepoID: 13, Tag: "Web"},
		{RepoID: 14, Tag: "Tools"},
	}
	_tags = append(_tags, excpected...)

	if err := TagRepository.SaveAndDeleteUnfindTags(
		context.Background(),
		_tags,
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer Repositories.Tag.DeleteMany(
		context.Background(),
		bson.M{"repo_id": bson.M{"$in": []uint{13,14,12} }},
		nil,
		options.Delete(),
	)

	tags, err := TagRepository.GetFilteredTagsByRepoID(
		context.Background(),
		12,
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !(tags[0].Tag == excpected[0].Tag && tags[1].Tag == excpected[1].Tag) {
		t.Log("Asserting error")
		t.FailNow()
	}
}