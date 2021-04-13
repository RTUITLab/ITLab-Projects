package repositories_test

import (
	"github.com/ITLab-Projects/pkg/githubreq"
	"testing"

	"github.com/ITLab-Projects/pkg/repositories"
)


var Repositories *repositories.Repositories

func init() {
	_r, err := repositories.New(&repositories.Config{
		DBURI: "mongodb://root:root@127.0.0.1:27100/ITLabProjects",
	})
	if err != nil {
		panic(err)
	}

	Repositories = _r
}

func TestFunc_Save(t *testing.T) {
	r := githubreq.New(&githubreq.Config{
		AccessToken: "",
	})
	repos, err := r.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if err := Repositories.Save(repos); err != nil {
		t.Log(err)
		t.FailNow()
	}
}