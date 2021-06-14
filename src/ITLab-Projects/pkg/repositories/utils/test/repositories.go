package test

import (
	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/pkg/repositories/utils"
	"os"
)

// Use only for test
// Validate uri
func GetTestRepository() *repositories.Repositories {
	dburi, find := os.LookupEnv("ITLAB_PROJECTS_DBURI")
	if !find {
		panic("Don't find dburi")
	}

	dburitest, find := os.LookupEnv("ITLAB_PROJECTS_DBURI_TEST")
	if !find {
		panic("Don't find dburi")
	}

	utils.ValidateTestURI(
		dburi,
		dburitest,
	)

	r, err := repositories.New(&repositories.Config{
		DBURI: dburitest,
	})
	if err != nil {
		panic(err)
	}

	return r
}