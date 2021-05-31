package v1

import (
	"regexp"

	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"

	"github.com/ITLab-Projects/service/api/v1/estimate"
	"github.com/ITLab-Projects/service/api/v1/functask"
	"github.com/ITLab-Projects/service/api/v1/issues"
	"github.com/ITLab-Projects/service/api/v1/projects"
	"github.com/ITLab-Projects/service/api/v1/tags"
	"github.com/gorilla/mux"
)

var regadmin *regexp.Regexp

func init() {
	regadmin = regexp.MustCompile(`(?m)(_|^)admin(_|$)`)
}

func (a *Api) buildAdmin(
	name	string,
	route	*mux.Route,
) {
	if regadmin.MatchString(name) {
		log.Debugf("Match on %s", name)
		handler := route.GetHandler()
		route.Handler(
			auth.AdminMiddleware(
				handler,
			),
		)
	}
}

func (a *Api) buildMongoSession(
	route *mux.Route,
) {
	route.Handler(
		mgsess.PutSessionINTOCtx(
			route.GetHandler(),
		),
	)
}

func (a *Api) BuildMiddlewares(
	route *mux.Route, 
	router *mux.Router, 
	ancestors []*mux.Route,
) error {
	a.buildMongoSession(route)
	
	name := route.GetName()
	if !a.Testmode {
		a.buildAdmin(
			name,
			route,
		)
	}

	return nil
}

func (a *Api) buildEndpoints() ServiceEndpoints {
	endpoints := ServiceEndpoints{
		Projects: projects.MakeEndpoints(a.projectService),
		Issues: issues.MakeEndPoints(a.issueService),
		Tags: tags.MakeEndpoints(a.tagsService),
		Task: functask.MakeEndPoints(a.taskService),
		Est: estimate.MakeEndPoints(a.estService),
	}

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

	return endpoints
}