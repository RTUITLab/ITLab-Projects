# ITLab-Projects
Service for listing RTUITLab projects work

REST API requests: https://www.postman.com/collections/a312d4a3f8da79bacc50
## Configuration

File ```src/api/auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json
    "audience": "example_audience", //audince for JWT
    "issuer" : "https://exampleissuersite.com", //issuer for JWT
    "scope" : "my_scope", //required scope for JWT
    "Github": {
          "accessToken" : "github_access_token"
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
    "host": "exampledb",
    "dbPort": "27017",
    "dbName" : "ITLabProjects",
    "projectsCollectionName" : "projects",
    "reposCollectionName" : "repos",
    "labelsCollectionName" : "labels",
    "issuesCollectionName" : "issues"
  },
  "AppOptions": {
    "appPort": "8080",
    "testMode": false,      //testMode=true disables jwt validation
    "projectFileBranch": "develop", //on which branch project_info.json is situated
    "elemsPerPage" : 40     //content is displayed with pagination
  }
}
```

## Installation using Docker
Install Docker and in ```src``` directory write this code:
```
docker-compose up -d
```
If you’re using Docker natively on Linux, Docker Desktop for Mac, or Docker Desktop for Windows, then the server will be running on
```http://localhost:8080```