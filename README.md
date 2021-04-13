# ITLab-Projects
Service for listing RTUITLab projects work

[![Build Status](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_apis/build/status/ITLab/ITLab-Projects?branchName=master)](https://dev.azure.com/rtuitlab/RTU%20IT%20Lab/_build/latest?definitionId=139&branchName=master)

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
    "uri": "mongodb://user:password@localhost:27017/ITLabProjects"    // env: ITLABPROJ_URI
  },
  "AppOptions": {
    "appPort": "8080", // env: ITLABPROJ_APPPORT
    "testMode": false,      //testMode=true disables jwt validation | env: ITLABPROJ_TESTMODE
    "projectFileBranch": "develop", //on which branch project_info.json is situated | env: ITLABPROJ_PROJFILEBRANCH
    "elemsPerPage" : 40     //content is displayed with pagination | env: ITLABPROJ_ELEMSPERPAGE
  }
}
```

## Build 
### Requirements
- Go 1.13.8+

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