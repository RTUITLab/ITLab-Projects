package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//TODO: implement gitlab projects paging (it's limited by 20 projects)


func getRepsFromGithub() []models.Repos{
	tempReps := make([]models.Repos, 0)
	URL := "https://api.github.com/orgs/RTUITLab/repos"

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

	json.NewDecoder(resp.Body).Decode(&tempReps)

	for i := range tempReps {
		tempReps[i].Platform = "github"
	}
	return tempReps
}

func getRepsFromGitlab() []models.Repos {
	tempReps := make([]models.Repos, 0)

	URL := "https://gitlab.com/api/v4/groups/6526027/projects?include_subgroups=true"

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
		return nil
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempReps)
	for i, repos := range tempReps {
		tempReps[i].Platform = "gitlab"
		tempReps[i].HTMLUrl = repos.GitLabHTMLUrl
		tempReps[i].UpdatedAt = repos.GitLabUpdatedAt
		tempReps[i].GitLabHTMLUrl = ""
		tempReps[i].GitLabUpdatedAt = ""
	}
	return tempReps
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