# ITLab-Projects
Service for listing RTUITLab projects work

REST API requests: https://www.postman.com/collections/a312d4a3f8da79bacc50
## Configuration

File ```src/ITLabReports/api/auth_config.json``` must contain next content:

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
File ```src/ITLabReports/api/config.json``` must contain next content:

```js
{
  "DbOptions": {
    "host": "exampledb",
    "dbPort": "27017",
    "dbName" : "ITLabProjects",
    "projectsCollectionName" : "projects",
    "reposCollectionName" : "repos"
  },
  "AppOptions": {
    "appPort": "8080",
    "testMode": false,      //testMode=true disables jwt validation
    "elemsPerPage" : 40     //content is displayed with pagination
  }
}
```