package githubreq_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/ITLab-Projects/pkg/githubreq"
)

var requster *githubreq.GHRequester

func init() {
	// Strange path but okay
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	token, find := os.LookupEnv("ITLAB_PROJECTS_ACCESSKEY")
	if !find {
		panic("Don't find token")
	}

	requster = githubreq.New(
		&githubreq.Config{
			AccessToken: token,
		},
	)
}

func TestFunc_GetRepositoris(t *testing.T) {
	repos, err := requster.GetRepositories()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("%v \n", repos)

	for _, r := range repos {
		t.Logf("name:%s langs: %v conts: %v\n", r.Name, r.Languages, r.Contributors)
	}
}

func TestFunc_URL(t *testing.T) {
	baseUrl := url.URL{
		Scheme: "https",
		Host: "api.github.com",
		Path: "orgs/RTUITLab",
	}

	val := url.Values{}
	val.Add("page", "10") 
	val.Add("count", "20")
	baseUrl.RawQuery = val.Encode()
	t.Log(baseUrl.String())
}