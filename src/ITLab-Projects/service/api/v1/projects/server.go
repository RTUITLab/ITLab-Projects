package projects

import (
	"context"

	"github.com/ITLab-Projects/service/api/v1/encoder"
	"github.com/ITLab-Projects/service/api/v1/errorencoder"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	UpdateProjectName	string = "update_projects_admin"
	GetProjectName		string = "get_project"
	GetProjectsName		string = "get_projects"
	DeleteProjectName	string = "delete_project_admin"
)

func NewHTTPServer(
	ctx			context.Context,
	endpoints	Endpoints,
	r			*mux.Router,
) {
	r.Handle(
		"/projects",
		httptransport.NewServer(
			endpoints.UpdateProjects,
			decodeUpdateProjetcsReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("POST").Name(UpdateProjectName)

	r.Handle(
		"/projects/{id:[0-9]+}",
		httptransport.NewServer(
			endpoints.GetProject,
			decodeGetProjectReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("GET").Name(GetProjectName)

	r.Handle(
		"/projects",
		httptransport.NewServer(
			endpoints.GetProjects,
			decodeGetProjectsReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("GET").Name(GetProjectsName)

	r.Handle(
		"/projects/{id:[0-9]+}",
		httptransport.NewServer(
			endpoints.DeleteProject,
			decodeDeleteProjectsReq,
			encoder.EncodeResponce,
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods("DELETE").Name(DeleteProjectName)
}