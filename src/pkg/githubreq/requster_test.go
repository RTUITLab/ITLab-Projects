package githubreq_test

import (
	"net/url"
	"testing"

	"github.com/ITLab-Projects/pkg/githubreq"
)

var requster *githubreq.GHRequester

func init() {
	requster = githubreq.New(
		&githubreq.Config{
			AccessToken: "",
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