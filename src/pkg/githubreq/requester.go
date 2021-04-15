package githubreq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/user"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

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
		},
		accessToken: cfg.AccessToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
		},
	}
}

const (
	scheme = "https"
	host = "api.github.com"
	orgName = "RTUITLab"
)

type GHRequester struct {
	client 			*http.Client

	baseUrl 		url.URL

	accessToken		string

	maxRepsPage		int
}


// GetRepositories return all repositories from GitHub
func (r *GHRequester) getRepositories(kv ...string) ([]repo.RepoWithURLS, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("orgs/%s/repos", orgName)
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
		return nil, errors.Wrapf(UnexpectedCode, "%v", resp.StatusCode)
	}

	var repos []repo.RepoWithURLS
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *GHRequester) GetMilestonesForRepo(repName string) ([]milestone.Milestone, error) {
	issues, err := r.getAllIssues(repName)
	if err != nil {
		return nil, err
	}

	return r.getAllMilestones(issues), nil
}

// return buffer with resp body
func (r *GHRequester) getAllIssues(repName string) ([]milestone.IssueFromGH, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("repos/%s/%s/issues", orgName, repName)
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
		return nil, errors.Wrapf(UnexpectedCode, "%v", resp.StatusCode)
	}

	var issues []milestone.IssueFromGH

	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func (r *GHRequester) getAllMilestones(issues []milestone.IssueFromGH) ([]milestone.Milestone) {
	set := make(map[interface{}][]milestone.Issue)

	for _, issue := range issues {
		if _, find := set[issue.Milestone]; issue.Milestone != nil && !find {
			set[issue.Milestone] = []milestone.Issue{issue.Issue}
		} else if issue.Milestone != nil && find {
			set[issue.Milestone] = append(set[issue.Milestone], issue.Issue)
		}
	}

	var milestones []milestone.Milestone

	for k, v := range set {
		m := k.(*milestone.MilestoneFromGH)
		milestones = append(milestones,  milestone.Milestone{MilestoneFromGH: *m, Issues: v})
	}

	return milestones
}

func (r *GHRequester) GetRepositoriesWithoutURL() ([]repo.Repo, error) {
	reps, err := r.GetRepositories()
	if err != nil {
		return nil, err
	}

	return repo.ToRepo(reps), nil
}

func (r *GHRequester) GetRepositories() ([]repo.RepoWithURLS, error) {
	if err := r.getMaxRepPage(); err != nil {
		return nil, err
	}

	// Get All Repos
	repoChan := make(chan []repo.RepoWithURLS, r.maxRepsPage)
	{
		var wg sync.WaitGroup
		for i := 1; i <= r.maxRepsPage; i++ {
			wg.Add(1)
			go func(i int, wg *sync.WaitGroup){
				defer wg.Done()
				reps, err := r.getRepositories("page", fmt.Sprint(i), "sort", "updated", "direction", "asc")
				if err != nil {
					log.WithFields(
						log.Fields{
							"package": "githubreq",
							"func": "GetRepositories",
							"err": err,
						},
					).Warn("Failed to get repo")
				} else {
					repoChan <- reps
				}

			}(i, &wg)
		}
		wg.Wait()
		close(repoChan)
	}

	var reps []repo.RepoWithURLS

	for rep := range repoChan {
		reps = append(reps, rep...)
	}

	// GetLangueages for each repo
	var wg sync.WaitGroup
	{
		setLanguages := func(rep *repo.RepoWithURLS, wg *sync.WaitGroup) {
			defer wg.Done()
			langs, err := r.getRepoLanguagesByURL(rep.LangaugesURL)
			if err != nil {
				log.WithFields(
					log.Fields{
						"package": "githubreq",
						"func": "GetRepositories",
						"err": err,
					},
				).Warnf("Failed to load language for repo: %s", rep.Name)
			} else {
				rep.Languages = langs
			}
		}

		setContributors := func(rep *repo.RepoWithURLS, wg *sync.WaitGroup) {
			defer wg.Done()
			users, err := r.getRepoContributorsByURL(rep.ContributorsURL)
			if err != nil {
				log.WithFields(
					log.Fields{
						"package": "githubreq",
						"func": "GetRepositories",
						"err": err,
					},
				).Warnf("Failed load contributors for repo %s", rep.Name)
			} else {
				rep.Contributors = users
			}
		}

		for i, _ := range reps {
			wg.Add(2)
			go setLanguages(&reps[i], &wg)
			go setContributors(&reps[i], &wg)
		}
	}
	wg.Wait()

	return reps, nil
}

func (r *GHRequester) getRepoLanguages(repName string) (map[string]int, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("repos/%s/%s/languages", orgName, repName)
	
	return r.getRepoLanguagesByURL(url.String())
}

func (r *GHRequester) getRepoLanguagesByURL(url string) (map[string]int, error) {
	var langs map[string]int
	req, err := http.NewRequest("GET", url, nil)
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
		return nil, errors.Wrapf(UnexpectedCode, "%v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&langs); err != nil {
		return nil, err
	}

	return langs, nil
}

func (r* GHRequester) getRepoContributorsByURL(url string) ([]user.User, error) {
	var users []user.User

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, errors.Wrapf(UnexpectedCode, "status code: %v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// prepareReqToGH set Authorization Header if them defined
func (r *GHRequester) prepareReqToGH(req *http.Request) {
	if r.accessToken != "" {
		// TODO read docs again and checkout of some new solution
		req.Header.Set("Authorization", "Bearer " + r.accessToken)
	}
}

// Resp should be from getting repos
func (r *GHRequester) setMaxRepPage(resp *http.Response) error {
	const pattern = `(?m)page=(\w+)>;\srel="last"`
	link := resp.Header.Get("Link")
	re := regexp.MustCompile(pattern)
	all := re.FindStringSubmatch(link)
	if len(all) != 2 {
		return errors.New("Can't get last page")
	}
	maxPage, _ := strconv.ParseInt(all[1], 10, 64)
	r.maxRepsPage = int(maxPage)
	return nil
}

func (r *GHRequester) getMaxRepPage() error {
	url := r.baseUrl
	url.Path += fmt.Sprintf("orgs/%s/repos", orgName)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	r.prepareReqToGH(req)
	
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := r.setMaxRepPage(resp); err != nil {
		return err
	}

	return nil
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