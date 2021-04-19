package apibuilder

import (
	"github.com/gorilla/mux"
)

type ApiBulder interface {
	Build(*mux.Router)
}

