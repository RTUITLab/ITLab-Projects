package server

import (
	"ITLab-Projects/models"
	"ITLab-Projects/server/utils"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getRepsFromGithub(page string, c chan models.Response) {
	tempReps := make([]models.Repos, 0)
	var wg sync.WaitGroup
	pageCount := 0
	all := false

	if page == "all" {
		all = true
		page = "1"
	}
	URL := "https://api.github.com/orgs/RTUITLab/repos?sort=updated&direction=asc&page=" + page
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
		wg.Add(2)
		go getRepLanguages(&tempReps[i], &wg)
		go getRepContributors(&tempReps[i], &wg)
		tempReps[i].Platform = "github"
		tempReps[i].Path = tempReps[i].Name
	}
	wg.Wait()
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

func getProjectInfoFile(rep *models.Repos, c chan models.ProjectInfo) {
	projectInfo := models.NewProjectInfo()
	fileUrl := "https://raw.githubusercontent.com/RTUITLab/" + rep.Path + "/" + cfg.App.ProjectFileBranch + "/project_info.json"
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
		rep.Meta = projectInfo.Repos
		projectInfo.Project.LastUpdated = rep.PushedAt
		for i := range projectInfo.Repos.Labels {
			projectInfo.Repos.Labels[i].Color = utils.MakeLabelColor(projectInfo.Repos.Labels[i].Name)
		}
	} else {
		projectInfo.Project.Path = rep.Path
		projectInfo.Project.LastUpdated = rep.PushedAt
		projectInfo.Project.Description = rep.Description
		projectInfo.Project.Reps = append(projectInfo.Project.Reps, rep.Path)

		projectInfo.Repos.Path = rep.Path
		projectInfo.Repos.HumanName = rep.Path
		projectInfo.Repos.Description = rep.Description
	}
	saveProjectToDB(projectInfo)
	getRepIssues(rep, projectInfo)
	c <- projectInfo
}

func saveProjectToDB(projectInfo models.ProjectInfo) {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"path": projectInfo.Project.Path}
	update := bson.M{
		"$set" : bson.M{
			"humanName" : projectInfo.Project.HumanName,
			"description" : projectInfo.Project.Description,
			"lastUpdated" : projectInfo.Project.LastUpdated,
		},
		"$addToSet" : bson.M{
			"reps" : projectInfo.Repos.Path,
			"labels" : bson.M{
				"$each" : projectInfo.Repos.Labels,
			},
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
		_, err := repsCollection.ReplaceOne(ctx, filter, rep, opts)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "mongodb.ReplaceOne",
				"handler":  "saveReposToDB",
				"rep":  rep.Path,
				"error":    err,
			},
			).Warn("Project update failed!")
		}
	}
}

func saveIssuesToDB(issues []models.Issue) {
	for _, issue := range issues {
		opts := options.Replace().SetUpsert(true)
		filter := bson.M{"id": issue.ID}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := issuesCollection.ReplaceOne(ctx, filter, issue, opts)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "mongodb.UpdateOne",
				"handler":  "saveIssuesToDB",
				"error":    err,
			},
			).Warn("Issues update failed!")
		}
	}
}

func saveIssueLabelsToDB(issues []models.Issue) {
	for _, issue := range issues {
		for _, label := range issue.Labels {
			label.Type = "rep"
			opts := options.Replace().SetUpsert(true)
			filter := bson.M{"name": label.Name}

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := labelsCollection.ReplaceOne(ctx, filter, label, opts)
			if err != nil {
				log.WithFields(log.Fields{
					"function": "mongodb.UpdateOne",
					"handler":  "saveIssueLabelsToDB",
					"error":    err,
				},
				).Warn("Issue Labels update failed!")
			}
		}
	}
}


func saveLabelsToDB(repos []models.Repos) {
	for _, rep := range repos {
		for _, label := range rep.Meta.Labels {
			opts := options.Replace().SetUpsert(true)
			filter := bson.M{"name": label.Name}

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := labelsCollection.ReplaceOne(ctx, filter, label, opts)
			if err != nil {
				log.WithFields(log.Fields{
					"function": "mongodb.UpdateOne",
					"handler":  "saveLabelsToDB",
					"rep":  rep.Path,
					"error":    err,
				},
				).Warn("Project update failed!")
			}
		}
	}
}

func saveLabelToDB(label models.Label) {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"name": label.Name}
	label.Type = "rep"

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := labelsCollection.ReplaceOne(ctx, filter, label, opts)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "mongodb.UpdateOne",
			"handler":  "saveLabelsToDB",
			"label":  label.Name,
			"error":    err,
		},
		).Warn("Label save failed!")
	}
}

func getRepLanguages(rep *models.Repos, wg *sync.WaitGroup) {
	var langs map[string]int
	URL := fmt.Sprintf("https://api.github.com/repos/RTUITLab/%s/languages", rep.Name)
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getLanugages",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&langs)
	rep.Languages = langs
	wg.Done()
}

func getRepContributors(rep *models.Repos, wg *sync.WaitGroup) {
	var contributors []models.User
	URL := fmt.Sprintf("https://api.github.com/repos/RTUITLab/%s/contributors", rep.Name)
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getLanugages",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&contributors)
	rep.Contributors = contributors
	wg.Done()
}

func getRepIssues(rep *models.Repos, projectInfo models.ProjectInfo) {
	var issues []models.Issue
	URL := fmt.Sprintf("https://api.github.com/repos/RTUITLab/%s/issues?state=all", rep.Name)
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer " + cfg.Auth.Github.AccessToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getRepIssues",
			"url"	: URL,
			"error"	:	err,
		},
		).Warn("Can't reach API!")
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&issues)
	for i := range issues {
		issues[i].RepPath = rep.Path
		issues[i].ProjectPath = projectInfo.Project.Path
		issues[i].Labels = append(issues[i].Labels, projectInfo.Repos.Labels...)
	}
	saveIssueLabelsToDB(issues)
	saveIssuesToDB(issues)
}

func calcPageTotal(repsTotal int64) int {
	pageTotal := int(repsTotal) / cfg.App.ElemsPerPage
	if int(repsTotal) % cfg.App.ElemsPerPage == 0 {
		return pageTotal
	} else {
		return pageTotal + 1
	}
}

func saveToDB(data interface{}) {
	var collection *mongo.Collection
	var dataSlice []interface{}

	switch e := data.(type) {
	case models.Issue, []models.Issue:
		collection = issuesCollection
	case models.Repos, []models.Repos:
		collection = repsCollection
	default:
		fmt.Printf("I don't know about type %T!\n", e)
		return
	}

	dataSlice = append(dataSlice, data)
	for _, elem := range dataSlice {
		opts := options.Replace().SetUpsert(true)
		filter := bson.M{"id": reflect.ValueOf(elem).FieldByName("ID").Uint()}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := collection.ReplaceOne(ctx, filter, elem, opts)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "mongodb.ReplaceOne",
				"handler": "saveToDB",
				"error": err,
			},
			).Warn("Data save failed!")
		}
	}

}