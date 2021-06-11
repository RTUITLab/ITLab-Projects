package githubreq

import (
	"github.com/ITLab-Projects/pkg/landingparser"
	"github.com/ITLab-Projects/pkg/models/landing"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/ITLab-Projects/pkg/clientwrapper"
	"github.com/ITLab-Projects/pkg/models/content"
	"github.com/ITLab-Projects/pkg/models/realese"
	"github.com/ITLab-Projects/pkg/models/tag"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/ITLab-Projects/pkg/models/user"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/pkg/models/repo"

	"net/url"
)

type Config struct {
	// AccessToken to GitGub

	AccessToken		string 
}

func New(cfg *Config) *GHRequester {
	r :=  &GHRequester {
		baseUrl: url.URL{
			Scheme: scheme,
			Host: host,
		},
		accessToken: cfg.AccessToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
				// DisableKeepAlives: true,
			},
		},
	}

	r.clientWithWrap = clientwrapper.New(r.client)
	r.clientWithWrap.Wrap(r.prepareReqToGH)

	return r
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

	clientWithWrap	*clientwrapper.ClientWithWrap
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

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

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

func (r *GHRequester) GetMilestonesForRepoWithID(rep repo.Repo) ([]milestone.MilestoneInRepo, error) {
	ms, err := r.GetMilestonesForRepo(rep.Name)
	if err != nil {
		return nil, err
	}

	var milestones []milestone.MilestoneInRepo

	for _, m := range ms {
		milestones = append(milestones, milestone.MilestoneInRepo{Milestone: m, RepoID: rep.ID})
	}

	return milestones, nil
}

func (r *GHRequester) GetAllMilestonesForRepoWithID(
	ctx context.Context,
	reps []repo.Repo,
	// error handling
	// if error is nil would'nt call
	f func(error),
) ([]milestone.MilestoneInRepo, error) {
	var ms []milestone.MilestoneInRepo
	msChan := make(chan []milestone.MilestoneInRepo)

	var count int = 0
	for i, _ := range reps {
		count++
		go func(rep repo.Repo){
			defer catchPanic()
			m, err := r.GetMilestonesForRepoWithID(rep)
			if err != nil {
				if f != nil {
					f(err)
				}
				msChan <- nil
				return
			}

			msChan <- m
		}(reps[i])
	}

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			close(msChan)
			return nil, ctx.Err()
		case m := <- msChan:
			ms = append(ms, m...)
		}
	}

	return ms, nil
}

func (r *GHRequester) GetAllTagsForRepoWithID(
	ctx context.Context,
	reps []repo.Repo,
	// if f nil would'nt call
	// if canceled return nil
	f func(error),
) ([]tag.Tag, error) {
	var tags []tag.Tag

	tgsChan := make(chan []tag.Tag)

	var count = 0
	for i, _ := range reps {
		count++
		go func(rep repo.Repo) {
			defer catchPanic()
			c, err := r.getLandingForRepo(rep)
			if err != nil {
				if f != nil {
					f(err)
				}
				tgsChan <- nil
				return
			}

			t, err := r.getTagsByURL(*c)
			if err != nil {
				if f != nil {
					f(err)
				}
				tgsChan <- nil
				return
			}

			tgsChan <- t
		}(reps[i])
	}

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			close(tgsChan)
			return nil, ctx.Err()
		case t := <- tgsChan:
			tags = append(tags, t...)
		}
	}

	return tags, nil
}

func (r *GHRequester) getTagsByURL(c content.Content) ([]tag.Tag, error) {

	req, err := http.NewRequest("GET", c.DownloadURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "%v for repo_id %v", resp.StatusCode, c.RepoID)
	}

	re := regexp.MustCompile(`(?m)^#\sTags[\s]*?\n(\*\s[\s]*\w+[\s]*?\n?)+$`)
	
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := re.FindString(bytes.NewBuffer(data).String())

	re = regexp.MustCompile(`(?m)\*[\s]*?(\w+)[\s]*?`)

	var tags []tag.Tag

	for _, match := range re.FindAllStringSubmatch(list, -1) {
		tags = append(tags, tag.Tag{RepoID: c.RepoID, Tag: match[1]})
	}

	return tags, nil
}

func (r *GHRequester) GetAllLandingsForRepoWithID(
	ctx 	context.Context,
	reps	[]repo.Repo,
	f		func(error),
) ([]*landing.Landing, error) {
	var ls []*landing.Landing

	lsChan := make(chan []*landing.Landing)

	var count = 0
	for i, _ := range reps {
		count++
		go func(rep repo.Repo) {
			defer catchPanicWithMsg(rep.Name)
			c, err := r.getLandingForRepo(rep)
			if err != nil {
				if f != nil {
					f(err)
				}
				lsChan <- nil
				return
			}

			if err = r.getLastCommit(
				rep, 
				c,
			); err != nil {
				if f != nil {
					f(err)
				}
				lsChan <- nil
				return
			}

			data, err := c.GetContent()
			if err != nil {
				if f != nil {
					f(err)
				}
				lsChan <- nil
				return
			}
			parser := landingparser.New()
			l := parser.Parse(
				landingparser.PrepareLandingToParse(data),
			)

			for i, img := range l.Image {
				l.Image[i] = c.GetURLForContent(img)
			}

			l.RepoId = c.RepoID
			l.Date.Time = c.GetDate()

			lsChan <- []*landing.Landing{l}
		}(reps[i])
	}

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			close(lsChan)
			return nil, ctx.Err()
		case l := <- lsChan:
			ls = append(ls, l...)
		}
	}
	return ls, nil
}

func (r *GHRequester) getLandingForRepo(
	rep repo.Repo,
) (*content.Content, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("/repos/%s/%s/contents/LANDING.md", orgName, rep.Name)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "%v for repo %s", resp.StatusCode, rep.Name)
	}

	var content content.Content
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, err
	}
	content.RepoID = rep.ID

	return &content, nil
}

func (r *GHRequester) getLastCommit(
	rep	repo.Repo,
	c	*content.Content,
) error {
	url := r.baseUrl
	url.Path += fmt.Sprintf("/repos/%s/%s/commits", orgName, rep.Name)
	q := url.Query()
	q.Add("per_page", "1")
	url.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(UnexpectedCode, "%v for repo %s", resp.StatusCode, rep.Name)
	}

	var commits []*content.Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return err
	}

	if len(commits) != 1 {
		return fmt.Errorf("Don't find commits")
	}

	c.Commit = commits[0]

	return nil
}

func (r *GHRequester) getIssueByPage(repName string, page int) ([]milestone.IssueFromGH, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("repos/%s/%s/issues", orgName, repName)
	q := url.Query()
	q.Add("state", "all")
	q.Add("milestone", "*")
	q.Add("page", fmt.Sprint(page))
	url.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "%v for repo %s", resp.StatusCode, repName)
	}

	var issues []milestone.IssueFromGH

	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}


func (r *GHRequester) getAllIssues(repName string) ([]milestone.IssueFromGH, error) {
	url := r.baseUrl
	url.Path += fmt.Sprintf("repos/%s/%s/issues", orgName, repName)
	q := url.Query()
	q.Add("state", "all")
	q.Add("milestone", "*")
	url.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}


	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "%v for repo %s", resp.StatusCode, repName)
	}

	var issues []milestone.IssueFromGH

	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	maxPage, err := getMaxPage(resp)
	if err != nil { // Means we can't get last page so it's only page
		return issues, nil
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	defer cancel()

	var count int
	isChan := make(chan []milestone.IssueFromGH, 1)
	errChan := make(chan error, 1)
	for i := 2; i <= maxPage; i++ {
		count++
		go func(i int) {
			defer catchPanic()
			iss, err := r.getIssueByPage(repName, i)
			if err == ErrForbiden || err == ErrUnatorizared {
				cancel()
				errChan <- err
				return
			} else if errors.Is(err, UnexpectedCode) {
				log.Debug(err)
				isChan <- nil
				return
			} else if err != nil {
				cancel()
				errChan <- err
				return
			}

			isChan <- iss
		}(i)
	}

	for i := 0; i < count; i++ {
		select {
		case <- ctx.Done():
			err := <- errChan
			close(errChan)
			close(isChan)
			return nil, err
		case iss := <-isChan:
			issues = append(issues, iss...)
		}
	}

	return issues, nil
	
}

func (r *GHRequester) getAllMilestones(issues []milestone.IssueFromGH) ([]milestone.Milestone) {
	set := make(map[interface{}][]milestone.Issue)

	for _, issue := range issues {
		if issue.Milestone != nil {
			if _, find := set[*issue.Milestone]; !find {
				set[*issue.Milestone] = []milestone.Issue{issue.Issue}
			} else if find {
				set[*issue.Milestone] = append(set[*issue.Milestone], issue.Issue)
			}
		}
	}

	var milestones []milestone.Milestone

	for k, v := range set {
		m := k.(milestone.MilestoneFromGH)
		milestones = append(milestones,  milestone.Milestone{MilestoneFromGH: m, Issues: v})
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

func (r *GHRequester) GetLastsRealeseWithRepoID(
	ctx context.Context,
	reps []repo.Repo,
	f func(error),
) ([]realese.RealeseInRepo, error) {
	var rls []realese.RealeseInRepo

	rlChan := make(chan realese.RealeseInRepo)

	var count = 0
	for i, _ := range reps {
		count++
		go func(rep repo.Repo) {
			rl, err := r.GetLastRealeseWithRepoID(rep)
			if err != nil && f != nil{
				f(err)
				rlChan <- realese.RealeseInRepo{RepoID: 0}
				return
			}
			rlChan <- rl
		}(reps[i])
	}

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case rl := <- rlChan:
			if rl.RepoID != 0 {
				rls = append(rls, rl)
			}
		}
	}

	return rls, nil
}

func (r *GHRequester) GetLastRealeseWithRepoID(rep repo.Repo) (realese.RealeseInRepo, error) {
	var real realese.RealeseInRepo

	url := r.baseUrl
	url.Path += fmt.Sprintf("/repos/%s/%s/releases/latest", orgName, rep.Name)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return real, err
	}

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return real, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return real, err
	}

	if resp.StatusCode != http.StatusOK {
		return real, errors.Wrapf(UnexpectedCode, "status code: %v for repo: %s", resp.StatusCode, rep.Name)
	}

	if err := json.NewDecoder(resp.Body).Decode(&real); err != nil {
		return real, err
	}

	real.RepoID = rep.ID

	return real, nil
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

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "%v for url %s", resp.StatusCode, url)
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

	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(UnexpectedCode, "status code: %v for url %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// prepareReqToGH set Authorization Header if them defined
func (r *GHRequester) prepareReqToGH(req *http.Request) {
	if r.accessToken != "" {
		req.Header.Set("Authorization", "Bearer " + r.accessToken)
	}
}

// Resp should be from getting repos
func (r *GHRequester) setMaxRepPage(resp *http.Response) error {
	maxPage, err := getMaxPage(resp)
	if err != nil {
		return err
	}

	r.maxRepsPage = int(maxPage)
	return nil
}

func getMaxPage(resp *http.Response) (int, error) {
	const pattern = `(?m)page=(\w+)>;\srel="last"`
	link := resp.Header.Get("Link")
	re := regexp.MustCompile(pattern)
	all := re.FindStringSubmatch(link)
	if len(all) != 2 {
		return 0, ErrGetLastPage
	}
	maxPage, _ := strconv.ParseInt(all[1], 10, 64)
	return int(maxPage), nil
}

func (r *GHRequester) getMaxRepPage() error {
	url := r.baseUrl
	url.Path += fmt.Sprintf("orgs/%s/repos", orgName)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	
	resp, err := r.clientWithWrap.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkStatusIfForbiddenOrUnathorizated(resp); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(UnexpectedCode, "status code: %v", resp.StatusCode)
	}

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

func checkStatusIfForbiddenOrUnathorizated(resp *http.Response) error {
	if resp.StatusCode == http.StatusForbidden {
		return ErrForbiden
	} else if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnatorizared
	}

	return nil
}

func catchPanic() {
	if r := recover(); r != nil {
		logrus.WithFields(
			log.Fields{
				"package": "requster",
				"err": r,
			},
		).Error("Catch panic without name")
	}
}

func catchPanicWithMsg(msg string) {
	if r := recover(); r != nil {
		logrus.WithFields(
			log.Fields{
				"package": "requster",
				"err": r,
				"msg": msg,
			},
		).Error("Catch panic with name")
	}
}