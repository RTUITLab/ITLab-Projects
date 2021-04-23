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
            "post": {
                "description": "make all request to github to update repositories, milestones\nIf don't get from gh some repos delete it in db",
                "summary": "Update all projects",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/err.Message"
                        }
                    },
                    "502": {
                        "description": "Bad Gateway",
                        "schema": {
                            "$ref": "#/definitions/err.Err"
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
                "summary": "add estimate to milestone",
                "parameters": [
                    {
                        "description": "estimate that you want to add",
                        "name": "estimate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/estimate.Estimate"
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
                    "404": {
                        "description": "estimate not found",
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
        "/api/v1/projects/task": {
            "post": {
                "description": "add func task to milestone\nif func task is exist for milesotne will replace it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "add func task to milestone",
                "parameters": [
                    {
                        "description": "function task that you want to add",
                        "name": "functask",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/functask.FuncTask"
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
                    "404": {
                        "description": "func task not found",
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
        }
    },
    "definitions": {
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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
