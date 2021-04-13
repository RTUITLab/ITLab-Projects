package githubreq

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ITLab-Projects/pkg/models/repo"

	"net/url"

	"github.com/ITLab-Projects/pkg/models/page"
)

type Config struct {
	// AccessToken to GitGub

	AccessToken		string 
}

// TODO доделать
type Requester interface {
	// key-value pairs for query
	GetRepositories(...string) ([]page.Page, error)

	GetRepositoriesForEach(func(page.Page)) error
	GetRepositoriesPage(uint) ([]repo.Repo, error)
	GetRepositoryByName(string) (repo.Repo, error)
}

// TODO return Requester
func New(cfg *Config) *GHRequester {
	return &GHRequester {
		baseUrl: url.URL{
			Scheme: scheme,
			Host: host,
			Path: basepath,
		},
		accessToken: cfg.AccessToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

const (
	scheme = "https"
	host = "api.github.com"
	basepath = "orgs/RTUITLab"
)

type GHRequester struct {
	client 			*http.Client

	baseUrl 		url.URL

	accessToken		string
}

// GetRepositories return repositories from GitHub
func (r *GHRequester) GetRepositories(kv ...string) ([]repo.Repo, error) {
	url := r.baseUrl
	url.Path += "/repos"

	query := parseKeyValue(kv...)
	url.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	r.prepareReqToGH(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, UnexpectedCode
	}
	var repos []repo.Repo

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// prepareReqToGH set Authorization Header if them defined
func (r *GHRequester) prepareReqToGH(req *http.Request) {
	if r.accessToken != "" {
		// TODO read docs again and checkout of some new solution
		req.Header.Set("Authorization", "Bearer " + r.accessToken)
	}
}

// parseKeyValue parse slice of string to url.Values
func parseKeyValue(kv ...string) url.Values {
	val := url.Values{}
	for i, v := range kv {
		if i % 2 == 1 {
			val.Add(kv[i-1], v)
		}
	}

	return val
}