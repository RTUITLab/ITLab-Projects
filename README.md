# ITLab-Projects
---
[REST API requests for postman](https://www.postman.com/collections/4b43b349d416cb99319e)

## Documantation 
Can be open directly by swagger. All swagger files located in ```src/ITLab-Projects/docs```.
Or if the service is running in default settings http://localhost:8080/api/projects/swagger/

## Configuration
File ```src/.env``` for laucnh from docker must contain:
```.env
# DATABASE Settings
MONGO_INITDB_ROOT_USERNAME=username
MONGO_INITDB_ROOT_PASSWORD=password
MONGO_INITDB_DATABASE=example_database

# App settings
ITLAB_PROJECTS_DBURI=mongodb://username:password@host:port/example_database
# database name for tests
ITLAB_PROJECTS_DBURI_TEST=mongodb://username:password@host:port/example_databaseTest
ITLAB_PROJECTS_ACCESSKEY=some_access_key
ITLABPROJ_ROLES="role role.admin"
# url to jwks.json
ITLABPROJ_KEYURL=https://example.com
# audince for JW
ITLABPROJ_AUDIENCE=audiance
# issuer for JWT
ITLABPROJ_ISSUER=https://example.com
ITLAB_PROJECTS_APPPORT=8080
# Test mode disable jwt validation and open /debug/pproh/ handlers to check app
ITLAB_PROJECTS_TESTMODE=false
# How many times will the projects be updated, 
# if it does not exist, it will not be updated itself
ITLAB_PROJECTS_UPDATETIME=2h
```
If you want run app directly from Go, copy this ```.env``` file to ```src/ITLab-Projects```

## Build
### Requirements
- Go 1.15+

in ```src/ITLab-Projects``` directory:
```
go build -o main
./main
```

## Build using docker

Install [Docker](https://www.docker.com) and in ```src``` write:
```
docker-compose up --build
```

If you’re using Docker natively on Linux, Docker Desktop for Mac, or Docker Desktop for Windows, then the server will be running on
```http://localhost:8080```