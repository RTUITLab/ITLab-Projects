# ITLab-Projects
Service for listing RTUITLab projects work

REST API requests: https://www.postman.com/collections/a312d4a3f8da79bacc50
## Configuration

File ```src/api/auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json | env: ITLABPROJ_KEYURL
    "audience": "example_audience", //audince for JWT | env: ITLABPROJ_AUDIENCE
    "issuer" : "https://exampleissuersite.com", //issuer for JWT | env: ITLABPROJ_ISSUER
    "scope" : "my_scope", //required scope for JWT | env: ITLABPROJ_SCOPE
    "Github": {
          "accessToken" : "github_access_token" // | env: ITLABPROJ_GHACCESSTOKEN
    },
    "Gitlab": {
          "accessToken" : "gitlab_access_token"
    }
  }
}
```

File ```src/api/config.json``` must contain next content:

```js
{
  "DbOptions": {
    "host": "exampledb",    // env: ITLABPROJ_HOST
    "dbPort": "27017", // env: ITLABPROJ_DBPORT
    "dbName" : "ITLabProjects", // env: ITLABPROJ_DBNAME
    "projectsCollectionName" : "projects", // env: ITLABPROJ_PROJCOLNAME
    "reposCollectionName" : "repos", // env: ITLABPROJ_REPSCOLNAME
    "labelsCollectionName" : "labels", // env: ITLABPROJ_LABSCOLNAME
    "issuesCollectionName" : "issues" // env: ITLABPROJ_ISSSCOLNAME
  },
  "AppOptions": {
    "appPort": "8080", // env: ITLABPROJ_APPPORT
    "testMode": false,      //testMode=true disables jwt validation | env: ITLABPROJ_TESTMODE
    "projectFileBranch": "develop", //on which branch project_info.json is situated | env: ITLABPROJ_ELEMSPERPAGE
    "elemsPerPage" : 40     //content is displayed with pagination | env: ITLABPROJ_PROJFILEBRANCH
  }
}
```

## Build 
### Requirements
- Go 1.12+

In ```src``` directory:
```
go build main.go
./main
```
## Build using Docker

Install Docker and in ```src``` directory write this code:
```
docker-compose up -d
```
If you’re using Docker natively on Linux, Docker Desktop for Mac, or Docker Desktop for Windows, then the server will be running on
```http://localhost:8080```