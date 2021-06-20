package apibuilder

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type ApiBulder interface {
	Build(*mux.Router)
	CreateServices()
	AddAuthMiddleware(endpoint.Middleware)
	AddLogger(log.Logger)
}

