package mgsess

import (
	"context"
	"encoding/json"
	"net/http"

	e "github.com/ITLab-Projects/pkg/err"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Kamva/mgm"
)

func PutSessionINTOCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logrus.Debug("Session middleware")
			_, client, _, _ := mgm.DefaultConfigs()
			sess, err := client.StartSession()
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
			defer sess.EndSession(
				context.Background(),
			)

			sessctx := mongo.NewSessionContext(
				r.Context(),
				sess,
			)


			r = r.WithContext(sessctx)

			next.ServeHTTP(w, r)
		},
	)
}