package typechecker_test

import (
	"reflect"

	"github.com/ITLab-Projects/pkg/models/repo"
	"github.com/ITLab-Projects/pkg/repositories/typechecker"

	"testing"
)

func TestFunc_NewTypeChecker(t *testing.T) {
	f := typechecker.NewSingle(reflect.TypeOf(repo.Repo{}))

	if err := f(&[]repo.Repo{}); err != nil {
		t.Log(err)
	}
}

func TestFunc_Test(t *testing.T) {
	repo := repo.Repo{ID: 10}

	v := reflect.ValueOf(&repo)
	field := v.FieldByName("ID")
	u := field.Uint()

	t.Log(u)
}
