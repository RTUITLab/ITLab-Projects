package server

import (
	"ITLab-Projects/models"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
)

var ReposList []models.Repos

var reposType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Repos",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"html_url": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
			"language": &graphql.Field{
				Type: graphql.String,
			},
			"archived": &graphql.Field{
				Type: graphql.Boolean,
			},
			"open_issues": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"repos": &graphql.Field{
				Type: reposType,
				Description: "One repository",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					nameQuery, isOK := params.Args["name"].(string)
					if isOK {
						for _, repos := range ReposList {
							if repos.Name == nameQuery {
								return repos, nil
							}
						}
					}
					return models.Repos{}, nil
				},
			},
			"reposList": &graphql.Field{
				Type: graphql.NewList(reposType),
				Description: "List of Repositories",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return ReposList, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func graphQL(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func init() {
	resp, err := http.Get("https://api.github.com/orgs/RTUITLab/repos")
	json.NewDecoder(resp.Body).Decode(&ReposList)
	if err != nil {
		fmt.Print("Error:", err)
	}
}
