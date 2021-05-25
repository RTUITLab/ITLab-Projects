package v1

import (
	"regexp"

	"github.com/ITLab-Projects/service/middleware/auth"
	"github.com/ITLab-Projects/service/middleware/mgsess"
	log "github.com/sirupsen/logrus"

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