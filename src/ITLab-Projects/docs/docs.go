// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/projects/": {
            "get": {
                "description": "return a projects you can filter count of them\ntags, name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "return projects according to query value",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "represents the number of skiped projects",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "represent a limit of projects",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "use to filter projects by tag",
                        "name": "tag",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "use to filter by name",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repoasproj.RepoAsProjCompact"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Failed to get repositories",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            },
            "post": {
                "description": "make all request to github to update repositories, milestones",
                "tags": [
                    "projects"
                ],
                "summary": "Update all projects",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/err.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/estimate": {
            "post": {
                "description": "add estimate to milestone\nif estimate is exist for milesotne will replace it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "add estimate to milestone",
                "parameters": [
                    {
                        "description": "estimate that you want to add",
                        "name": "estimate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/estimate.EstimateFile"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Unexpected body",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "Don't find milestone with this id",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Failed to save estimate",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/estimate/{milestone_id}": {
            "delete": {
                "description": "delete estimate from database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "delete estimate from database",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "should be uint",
                        "name": "milestone_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "estimate not found",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "409": {
                        "description": "some problems with microfileservice",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Failed to delete estimate",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/issues": {
            "get": {
                "description": "return issues according to query params",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "issues"
                ],
                "summary": "return issues",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "represent how mush skip first issues",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "set limit of getting issues",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search to name of issues, title of milestones and repository names",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search of label name of issues",
                        "name": "tag",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/milestone.IssuesWithMilestoneID"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/issues/labels": {
            "get": {
                "description": "return all unique labels of issues",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "issues"
                ],
                "summary": "return labels",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/tags": {
            "get": {
                "description": "return all tags",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "return all tags",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/tag.Tag"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/task": {
            "post": {
                "description": "add func task to milestone\nif func task is exist for milesotne will replace it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "functask"
                ],
                "summary": "add func task to milestone",
                "parameters": [
                    {
                        "description": "function task that you want to add",
                        "name": "functask",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/functask.FuncTaskFile"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Unexpected body",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "Don't find milestone with this id",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Failed to save functask",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/task/{milestone_id}": {
            "delete": {
                "description": "delete functask from database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "functask"
                ],
                "summary": "delete functask from database",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "should be uint",
                        "name": "milestone_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "func task not found",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "409": {
                        "description": "some problems with microfileservice",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Failed to delete func task",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/{id}": {
            "get": {
                "description": "return a project according to id value in path",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "return project according to id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "a uint value of repository id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/repoasproj.RepoAsProj"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete project by id and all milestones in it",
                "tags": [
                    "projects"
                ],
                "summary": "delete project by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id of project",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "403": {
                        "description": "if you are nor admin",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "409": {
                        "description": "some problems with microfileservice",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "assignee.Assignee": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "html_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                }
            }
        },
        "err.Err": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "err.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "estimate.Estimate": {
            "type": "object",
            "properties": {
                "estimate_url": {
                    "type": "string"
                },
                "milestone_id": {
                    "type": "integer"
                }
            }
        },
        "estimate.EstimateFile": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "milestone_id": {
                    "type": "integer"
                }
            }
        },
        "functask.FuncTask": {
            "type": "object",
            "properties": {
                "func_task_url": {
                    "type": "string"
                },
                "milestone_id": {
                    "type": "integer"
                }
            }
        },
        "functask.FuncTaskFile": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "milestone_id": {
                    "type": "integer"
                }
            }
        },
        "label.Label": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "node_id": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "milestone.Issue": {
            "type": "object",
            "properties": {
                "assignees": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/assignee.Assignee"
                    }
                },
                "body": {
                    "type": "string"
                },
                "closed_at": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "html_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/label.Label"
                    }
                },
                "number": {
                    "type": "integer"
                },
                "project_path": {
                    "type": "string"
                },
                "pull_request": {
                    "$ref": "#/definitions/pullrequest.PullRequest"
                },
                "reppath": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/user.User"
                }
            }
        },
        "milestone.IssuesWithMilestoneID": {
            "type": "object",
            "properties": {
                "assignees": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/assignee.Assignee"
                    }
                },
                "body": {
                    "type": "string"
                },
                "closed_at": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted": {
                    "type": "boolean"
                },
                "html_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/label.Label"
                    }
                },
                "milestone_id": {
                    "type": "integer"
                },
                "number": {
                    "type": "integer"
                },
                "project_path": {
                    "type": "string"
                },
                "pull_request": {
                    "$ref": "#/definitions/pullrequest.PullRequest"
                },
                "repo_id": {
                    "type": "integer"
                },
                "reppath": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/user.User"
                }
            }
        },
        "milestone.Milestone": {
            "type": "object",
            "properties": {
                "closed_issues": {
                    "type": "integer"
                },
                "creator": {
                    "$ref": "#/definitions/user.User"
                },
                "deleted": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "estimate": {
                    "$ref": "#/definitions/estimate.Estimate"
                },
                "func_task": {
                    "$ref": "#/definitions/functask.FuncTask"
                },
                "id": {
                    "type": "integer"
                },
                "issues": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/milestone.Issue"
                    }
                },
                "number": {
                    "type": "integer"
                },
                "open_issues": {
                    "type": "integer"
                },
                "state": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "pullrequest.PullRequest": {
            "type": "object",
            "properties": {
                "diff_url": {
                    "type": "string"
                },
                "html_url": {
                    "type": "string"
                },
                "patch_url": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "realese.Realese": {
            "type": "object",
            "properties": {
                "html_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "repo.Repo": {
            "type": "object",
            "properties": {
                "archived": {
                    "type": "boolean"
                },
                "contributors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.User"
                    }
                },
                "created_at": {
                    "type": "string"
                },
                "deleted": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "html_url": {
                    "description": "Path\t\t\t\tstring\t\t\t` + "`" + `json:\"path_with_namespace,omitempty\"` + "`" + `",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "languages": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "pushed_at": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "repoasproj.RepoAsProj": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "number"
                },
                "last_realese": {
                    "$ref": "#/definitions/realese.Realese"
                },
                "milestones": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/milestone.Milestone"
                    }
                },
                "repo": {
                    "$ref": "#/definitions/repo.Repo"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tag.Tag"
                    }
                }
            }
        },
        "repoasproj.RepoAsProjCompact": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "number"
                },
                "repo": {
                    "$ref": "#/definitions/repo.Repo"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tag.Tag"
                    }
                }
            }
        },
        "tag.Tag": {
            "type": "object",
            "properties": {
                "tag": {
                    "type": "string"
                }
            }
        },
        "user.User": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "html_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "ITLab-Projects API",
	Description: "This is a server to get projects from github",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
