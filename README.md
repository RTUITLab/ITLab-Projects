# ITLab-Projects
Service for listing RTUITLab projects work

REST API requests: https://www.getpostman.com/collections/a312d4a3f8da79bacc50
## Configuration

File ```src/ITLabReports/api/auth_config.json``` must contain next content:

```js
{
  "AuthOptions": {
    "keyUrl": "https://examplesite/files/jwks.json", //url to jwks.json
    "audience": "example_audience", //audince for JWT
    "issuer" : "https://exampleissuersite.com", //issuer for JWT
    "scope" : "my_scope" //required scope for JWT
  }
}

