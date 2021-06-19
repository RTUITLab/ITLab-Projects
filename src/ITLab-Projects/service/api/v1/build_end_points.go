package v1

import (
	"github.com/ITLab-Projects/service/api/v1/updater"
	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	"github.com/go-kit/kit/endpoint"


	"github.com/ITLab-Projects/service/api/v1/estimate"
	"github.com/ITLab-Projects/service/api/v1/functask"
	"github.com/ITLab-Projects/service/api/v1/issues"
	"github.com/ITLab-Projects/service/api/v1/landing"
	"github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/ITLab-Projects/service/api/v1/tags"
)

func (a *Api) buildEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	// ---------- Estimate ----------
	endpoints.Est.AddEstimate = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Est.AddEstimate)

	endpoints.Est.DeleteEstimate = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Est.DeleteEstimate)
	// ----------		----------

	// ---------- Task ----------
	endpoints.Task.AddFuncTask = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Task.AddFuncTask)

	endpoints.Task.DeleteFuncTask = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Task.DeleteFuncTask)
	// ----------		----------

	// ---------- Tags ----------
	endpoints.Tags.GetAllTags = endpoint.Chain(
		a.NewAuth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Tags.GetAllTags)
	// ----------		----------

	// ---------- Issues ---------- 
	endpoints.Issues.GetIssues = endpoint.Chain(
		a.NewAuth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Issues.GetIssues)

	endpoints.Issues.GetLabels = endpoint.Chain(
		a.NewAuth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Issues.GetLabels)
	// ----------		----------

	// ---------- Projects ----------
	endpoints.Projects.DeleteProject = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.DeleteProject)

	endpoints.Projects.GetProject = endpoint.Chain(
		a.NewAuth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProject)

	endpoints.Projects.GetProjects = endpoint.Chain(
		a.NewAuth,
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProjects)
	// ----------		----------

	// ---------- Landing ----------
	endpoints.Landing.GetLanding = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Landing.GetLanding)

	endpoints.Landing.GetAllLandings = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Landing.GetAllLandings)
	// ----------		----------

	// ---------- Updater ----------

	endpoints.Update.UpdateProjects = endpoint.Chain(
		a.NewAuth,
		auth.EndpointAdminMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Update.UpdateProjects)
	// ----------		----------
	
	return endpoints
}

func (a *Api) _buildEndpoint() ApiEndpoints {
	endpoints := a.endpoints()

	// ---------- Estimate ----------
	endpoints.Est.AddEstimate = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Est.AddEstimate)

	endpoints.Est.DeleteEstimate = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Est.DeleteEstimate)
	// ----------		----------

	// ---------- Task ----------
	endpoints.Task.AddFuncTask = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Task.AddFuncTask)

	endpoints.Task.DeleteFuncTask = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Task.DeleteFuncTask)
	// ----------		----------

	// ---------- Tags ----------
	endpoints.Tags.GetAllTags = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Tags.GetAllTags)
	// ----------		----------

	// ---------- Issues ---------- 
	endpoints.Issues.GetIssues = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Issues.GetIssues)

	endpoints.Issues.GetLabels = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Issues.GetLabels)
	// ----------		----------

	// ---------- Projects ----------
	endpoints.Projects.DeleteProject = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.DeleteProject)

	endpoints.Projects.GetProject = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProject)

	endpoints.Projects.GetProjects = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Projects.GetProjects)
	// ----------		----------


	// ---------- Landing ----------
	endpoints.Landing.GetLanding = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Landing.GetLanding)

	endpoints.Landing.GetAllLandings = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Landing.GetAllLandings)
	// ----------		----------

	// ---------- Updater ----------

	endpoints.Update.UpdateProjects = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Update.UpdateProjects)
	// ----------		----------


	return endpoints
}

func (a *Api) endpoints() ApiEndpoints {
	return ApiEndpoints{
		Projects: projects.MakeEndpoints(a.projectService),
		Issues: issues.MakeEndPoints(a.issueService),
		Tags: tags.MakeEndpoints(a.tagsService),
		Task: functask.MakeEndPoints(a.taskService),
		Est: estimate.MakeEndPoints(a.estService),
		Landing: landing.MakeEndpoints(a.landingService),
		Update: updater.MakeEndpoints(a.updaterService),
	}
}