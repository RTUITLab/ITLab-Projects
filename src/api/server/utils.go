package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//TODO: implement gitlab projects paging (it's limited by 20 projects)



func getRepsFromGithub(page string) ([]models.Repos, int) {
	tempReps := make([]models.Repos, 0)

	URL := "https://api.github.com/orgs/RTUITLab/repos?page=" + page

	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepsFrom",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return nil, 0
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempReps)

	for i := range tempReps {
		tempReps[i].Platform = "github"
		tempReps[i].Path = tempReps[i].Name
	}

	linkHeader := resp.Header.Get("Link")
	lastPage := strings.LastIndex(linkHeader, "page=")
	linkHeader = linkHeader[lastPage:]
	linkHeader = strings.TrimLeft(linkHeader, "page=")
	linkHeader = strings.TrimRight(linkHeader, ">; rel=\"last\"")
	pageCount, _ := strconv.Atoi(linkHeader)
	return tempReps, pageCount
}

func getRepsFromGitlab(page string) ([]models.Repos, int) {
	tempReps := make([]models.Repos, 0)

	URL := "https://gitlab.com/api/v4/groups/6526027/projects?include_subgroups=true&page="+page

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepsFrom",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return nil, 0
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempReps)
	for i := range tempReps {
		tempReps[i].Platform = "gitlab"
		tempReps[i].Path = url.QueryEscape(tempReps[i].Path)
		tempReps[i].HTMLUrl = tempReps[i].GitLabHTMLUrl
		tempReps[i].UpdatedAt = tempReps[i].GitLabUpdatedAt
		tempReps[i].GitLabHTMLUrl = ""
		tempReps[i].GitLabUpdatedAt = ""
	}

	pageCountHeader := resp.Header.Get("X-Total-Pages")
	pageCount, _ := strconv.Atoi(pageCountHeader)
	return tempReps, pageCount
}

func getRepFromGithub(id string) models.Repos {
	var tempRep models.Repos
	URL := "https://api.github.com/repos/RTUITLab/" + id

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepFromGithub",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return models.Repos{}
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempRep)

	tempRep.Platform = "github"
	tempRep.Path = tempRep.Name

	return tempRep
}

func getRepFromGitlab(id string) models.Repos {
	var tempRep models.Repos
	URL := "https://gitlab.com/api/v4/projects/" + id

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepsFrom",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return models.Repos{}
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempRep)
	tempRep.Platform = "gitlab"
	tempRep.Path = url.QueryEscape(tempRep.Path)
	tempRep.HTMLUrl = tempRep.GitLabHTMLUrl
	tempRep.UpdatedAt = tempRep.GitLabUpdatedAt
	tempRep.GitLabHTMLUrl = ""
	tempRep.GitLabUpdatedAt = ""
	return tempRep
}

func getIssuesForGithub(id string, state string) []models.Issue {
	tempIssues := make([]models.Issue, 0)
	var URL string

	switch state {
	case "opened":
		URL = "https://api.github.com/repos/RTUITLab/"+id+"/issues?state=opened"
	case "closed":
		URL = "https://api.github.com/repos/RTUITLab/"+id+"/issues?state=closed"
	case "all":
		URL = "https://api.github.com/repos/RTUITLab/"+id+"/issues?state=all"
	default:
		URL = "https://api.github.com/repos/RTUITLab/"+id+"/issues"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepsFrom",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return nil
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&tempIssues)

	return tempIssues
}

func getIssuesForGitlab(id string, state string) []models.Issue {
	tempIssues := make([]models.Issue, 0)
	var URL string

	switch state {
	case "opened":
		URL = "https://gitlab.com/api/v4/projects/"+id+"/issues?state=opened"
	case "closed":
		URL = "https://gitlab.com/api/v4/projects/"+id+"/issues?state=closed"
	case "all":
		URL = "https://gitlab.com/api/v4/projects/"+id+"/issues?state=all"
	default:
		URL = "https://gitlab.com/api/v4/projects/"+id+"/issues"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getIssuesForGitlab",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return nil
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&tempIssues)

	for i := range tempIssues {
		tempIssues[i].User.ID = tempIssues[i].GitlabUser.ID
		tempIssues[i].User.Login = tempIssues[i].GitlabUser.GitLabLogin
		tempIssues[i].User.AvatarURL = tempIssues[i].GitlabUser.AvatarURL
		tempIssues[i].User.URL = tempIssues[i].GitlabUser.GitLabHTMLUrl
		tempIssues[i].GitlabUser = nil

		tempIssues[i].Number = *tempIssues[i].GitLabNumber
		tempIssues[i].GitLabNumber =  nil

		tempIssues[i].HtmlUrl = tempIssues[i].GitLabHTMLUrl
		tempIssues[i].GitLabHTMLUrl = ""
	}
	return tempIssues
}

func getIssueFromGithub(id string, number string) models.Issue {
	var tempIssue models.Issue
	URL := "https://api.github.com/repos/RTUITLab/"+id+"/issues/"+number

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepFromGithub",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return models.Issue{}
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&tempIssue)

	return tempIssue
}

func getIssueFromGitlab(id string, number string) models.Issue {
	var tempIssue models.Issue
	URL := "https://gitlab.com/api/v4/projects/"+id+"/issues/"+number

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepFromGithub",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		return models.Issue{}
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&tempIssue)

	tempIssue.User.ID = tempIssue.GitlabUser.ID
	tempIssue.User.Login = tempIssue.GitlabUser.GitLabLogin
	tempIssue.User.AvatarURL = tempIssue.GitlabUser.AvatarURL
	tempIssue.User.URL = tempIssue.GitlabUser.GitLabHTMLUrl
	tempIssue.GitlabUser = nil

	tempIssue.Number = *tempIssue.GitLabNumber
	tempIssue.GitLabNumber =  nil

	tempIssue.HtmlUrl = tempIssue.GitLabHTMLUrl
	tempIssue.GitLabHTMLUrl = ""
	
	return tempIssue
}