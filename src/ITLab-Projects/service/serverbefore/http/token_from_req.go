package http

import (
	log "github.com/sirupsen/logrus"
	ctxtoken "github.com/ITLab-Projects/pkg/conextvalue/token"
	"context"
	"net/http"
)

func TokenFromReq(
	ctx	context.Context,
	r	*http.Request,
) context.Context {
	log.WithFields(
		log.Fields{
			"package": "serverbefore/http",
			"func": "TokenFromReq",
		},
	).Debug("Token before")
	token := r.Header.Get("Authorization")
	
	return ctxtoken.New(
		ctx,
		token,
	)
}