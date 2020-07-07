package server

import (
	"ITLab-Projects/models"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getRepsFromGithub(page string, c chan models.Response) {
	var URL string
	tempReps := make([]models.Repos, 0)
	pageCount := 0
	all := false

	if page == "all" {
		all = true
		page = "1"
	}
	URL = "https://api.github.com/orgs/RTUITLab/repos?page=" + page

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
		c <- models.Response{}
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tempReps)

	for i := range tempReps {
		tempReps[i].Platform = "github"
		tempReps[i].Path = tempReps[i].Name
	}
	if len(tempReps) != 0 {
		linkHeader := resp.Header.Get("Link")
		lastPage := strings.LastIndex(linkHeader, "page=")
		linkHeader = linkHeader[lastPage:]
		linkHeader = strings.TrimLeft(linkHeader, "page=")
		linkHeader = strings.TrimRight(linkHeader, ">; rel=\"last\"")
		pageCount, _ = strconv.Atoi(linkHeader)
	}
	if all {
		cGithub := make(chan models.Response)
		for i:=2; i<=pageCount;i++ {
			go getRepsFromGithub(strconv.Itoa(i), cGithub)
		}
		for i:=2; i<=pageCount;i++ {
			repPage := <-cGithub
			tempReps = append(tempReps, repPage.Repositories...)
		}
	}
	response := models.Response{tempReps, pageCount}
	c <- response
}

func getRepsFromGitlab(page string, c chan models.Response) {
	tempReps := make([]models.Repos, 0)
	pageCount := 0

	URL := "https://gitlab.com/api/v4/groups/6526027/projects?include_subgroups=true&page="+page

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	req.Header.Set("Connection", "keep-alive")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepsFrom",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
		c <- models.Response{}
		return
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
	pageCount, err = strconv.Atoi(pageCountHeader)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "strconv.Atoi",
			"handler" : "getRepsFromGitlab",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't get pages count!")
	}

	response := models.Response{tempReps, pageCount}
	c <- response
}

func getRepFromGithub(id string) models.Repos {
	var tempRep models.Repos
	URL := "https://api.github.com/repos/RTUITLab/" + id

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
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

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := httpClient.Do(req)
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

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := httpClient.Do(req)
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

		tempIssues[i].Description = tempIssues[i].GitlabDescription
		tempIssues[i].GitlabDescription = ""

	}
	return tempIssues
}

func getIssueFromGithub(id string, number string) models.Issue {
	var tempIssue models.Issue
	URL := "https://api.github.com/repos/RTUITLab/"+id+"/issues/"+number

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
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

	
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Gitlab.AccessToken)
	resp, err := httpClient.Do(req)
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

	tempIssue.Description = tempIssue.GitlabDescription
	tempIssue.GitlabDescription = ""

	return tempIssue
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(5) * time.Second,
	}
	return client
}

func getProjectInfoFile(repPath string, c chan models.ProjectInfo) {
	var projectInfo models.ProjectInfo
	fileUrl := "https://raw.githubusercontent.com/RTUITLab/" + repPath + "/develop/project_info.json"

	req, err := http.NewRequest("GET", fileUrl, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "http.Get",
			"handler":  "getProjectInfoFile",
			"url":      fileUrl,
			"error":    err,
		},
		).Warn("Something went wrong")
		c <- models.ProjectInfo{}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&projectInfo)
		saveProjectToDB(projectInfo)
		c <- projectInfo
	} else {
		c <- models.ProjectInfo{}
	}
}

func saveProjectToDB(projectInfo models.ProjectInfo) {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"path": projectInfo.Project.Path}
	update := bson.M{
		"$set" : bson.M{
			"humanName" : projectInfo.Project.HumanName,
			"description" : projectInfo.Project.Description,
	},
		"$addToSet" : bson.M{
			"reps" : projectInfo.Repos.Path,
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := projectsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "mongodb.UpdateOne",
			"handler":  "saveProjectToDB",
			"project":  projectInfo.Project.Path,
			"error":    err,
		},
		).Warn("Project update failed!")
	}
}

func saveReposToDB(repos []models.Repos) {
	for _, rep := range repos {
		opts := options.Replace().SetUpsert(true)
		filter := bson.M{"id": rep.ID}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := reposCollection.ReplaceOne(ctx, filter, rep, opts)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "mongodb.UpdateOne",
				"handler":  "saveReposToDB",
				"rep":  rep.Path,
				"error":    err,
			},
			).Warn("Project update failed!")
		}
	}
}