package mgsess

import (
	"github.com/ITLab-Projects/pkg/repositories"
	"context"
	"encoding/json"
	"net/http"

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