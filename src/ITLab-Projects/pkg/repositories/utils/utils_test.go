package utils_test

import (
	"testing"

	"github.com/ITLab-Projects/pkg/repositories/utils"
)

func TestFunc_ValidareURI_Okay(t *testing.T) {
	utils.ValidateTestURI(
		"mongodb://user:password@net:27100/ITLabProjects?authSource=admin",
		"mongodb://user:password@net:27100/ITLabProjectsTest?authSource=admin",
	)
}

func TestFunc_ValidateURI_Panic_Equals(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("Assert error should be error")
			t.FailNow()
		}
	}()
	utils.ValidateTestURI(
		"mongodb://user:password@net:27100/ITLabProjects?authSource=admin",
		"mongodb://user:password@net:27100/ITLabProjects?authSource=admin",
	)
}

func TestFunc_ValidateURI_Panic_DontHaveTest(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("Assert error should be error")
			t.FailNow()
		}
	}()
	utils.ValidateTestURI(
		"mongodb://user:password@net:27100/ITLabProjects?authSource=admin",
		"mongodb://user:password@net:27100/itlab-projects?authSource=admin",
	)
}
