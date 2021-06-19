package v2

import (
	"github.com/ITLab-Projects/service/api/v2/projects"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	"github.com/go-kit/kit/endpoint"
)

func (a *Api) buildEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	// ---------- Projects ----------
	endpoints.Projects.GetProjects = endpoint.Chain(
		a.Auth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProjects)
	// ----------		----------

	return endpoints
}

func (a *Api) _buildEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	// ---------- Projects ----------
	endpoints.Projects.GetProjects = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProjects)
	// ----------		----------

	return endpoints
}

func (a *Api) endpoints() ApiEndpoints {
	return ApiEndpoints{
		Projects: projects.MakeEndpoints(a.projectService),
	}
}