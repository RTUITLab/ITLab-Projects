package mgsess

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ITLab-Projects/pkg/repositories"
	"github.com/ITLab-Projects/pkg/statuscode"
	"github.com/go-kit/kit/endpoint"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/sirupsen/logrus"
)

func PutSessionINTOCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logrus.Debug("Session middleware")
			sessctx, err := repositories.GetMongoSessionContext(
				r.Context(),
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(
					e.Message{
						Message: "Failed to connect to database",
					},
				)
				logrus.WithFields(
					logrus.Fields{
						"package": "middleware",
						"handler": "PutSessionINTOCtx",
						"err": err,
					},
				).Error()
				return
			}
			defer sessctx.EndSession(
				context.Background(),
			)

			r = r.WithContext(sessctx)

			next.ServeHTTP(w, r)
		},
	)
}

func PutMongoSessIntoCtx() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			logrus.Debug("Session middleware")
			sessctx, err := repositories.GetMongoSessionContext(
				ctx,
			)
			if err != nil {
				logrus.WithFields(
					logrus.Fields{
						"pkg": "middlewares/mgsess",
						"func": "PutMongoSessIntoCtx",
						"err": err,
					},
				).Error()
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("faield to get mongosession"),
					http.StatusInternalServerError,
				)
			}

			defer sessctx.EndSession(
				context.Background(),
			)

			return next(sessctx, request)
		}
	}
}