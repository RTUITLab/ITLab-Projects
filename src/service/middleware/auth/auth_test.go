package auth_test

import (
	"strings"
	"testing"

	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var rolesSet map[string]struct{}

func init() {
	rolesSet = map[string]struct{}{}
	var config config.Config
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}
	if err := envconfig.Process("itlab_projects", &config); err != nil {
		panic(err)
	}

	for _, r := range strings.Split(config.Auth.Roles, " ") {
		rolesSet[r] = struct{}{}
	}

}

func TestFunc_FindRole_Admin(t *testing.T) {
	role, err := auth.NewRoleGetter("itlab", rolesSet)(
		map[string]interface{} {
			"itlab": []string {
				"user",
    			"projects.admin",
			},
		},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(role)
}

func TestFunc_FindRole_User(t *testing.T) {
	role, err := auth.NewRoleGetter("itlab", rolesSet)(
		map[string]interface{} {
			"itlab": []string {
				"user",
    			"projects",
			},
		},
	)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(role)
}

func TestFunc_FindRole_DontFind(t *testing.T) {
	_, err := auth.NewRoleGetter("itlab", rolesSet)(
		map[string]interface{} {
			"itlab": []string {
				"user",
			},
		},
	)
	if err == nil {
		t.FailNow()
	}

	t.Log(err)
}

func TestFunc_FailedToCast(t *testing.T) {
	_, err := auth.NewRoleGetter("itlab", rolesSet)(
		map[string]interface{} {
			"itlab": nil,
		},
	)
	if err == nil {
		t.FailNow()
	}

	t.Log(err)
}