package auth

import (
	"regexp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ITLab-Projects/pkg/conextvalue"

	"github.com/ITLab-Projects/pkg/config"
	e "github.com/ITLab-Projects/pkg/err"
	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
)

type AuthMiddleware mux.MiddlewareFunc

func New(
	cfg config.AuthConfig,
) AuthMiddleware {
	rolesSet := map[string]struct{}{}

	for _, role := range strings.Split(cfg.RolesConfig.Roles, " ") {
		rolesSet[role] = struct{}{}
	}

	auth := AuthMiddleware(
		newAuthMiddleware(
			&cfg,
			NewRoleGetter("itlab", rolesSet),
		),
	)

	return auth
}

func newAuthMiddleware(
	cfg *config.AuthConfig,
	f getRoleFromClaim,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: cfg.KeyURL}, nil)
				configuration := auth0.NewConfiguration(client, []string{cfg.Audience}, cfg.Issuer, jose.RS256)
				validator := auth0.NewValidator(configuration, nil)
				token, err := validator.ValidateRequest(r)
				if err != nil {
					log.WithFields(log.Fields{
						"requiredAlgorithm" : "RS256",
						"error" : err,
					}).Warning("Token is not valid!")

					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(
						e.Message{
							Message: "Token is not valid",
						},
					)
					return
				}

				claims := map[string]interface{}{}
				if err = validator.Claims(r, token, &claims); err != nil {
					log.WithFields(log.Fields{
						"requiredClaims" : "iss, aud, sub, role",
						"error" : err,
					}).Warning("Invalid claims!")
		
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(
						e.Message {
							Message: "Invalid claims",
						},
					)
					return
				}
				
				role, err := f(claims)
				if err != nil {
					log.WithFields(log.Fields{
						"package" : "middleware/auth",
						"func": "authMiddleware",
						"error" : err,
					}).Warning("Failed to get role")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(
						e.Message{
							Message: "Failed to get role",
						},
					)
					return
				}

				ctx := context.WithValue(
					r.Context(),
					conextvalue.Role,
					role,
				)

				r = r.WithContext(ctx)
				
				next.ServeHTTP(w, r)
			},
		)
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			role := GetRoleFromCTX(r.Context())

			re := regexp.MustCompile(`\w+.admin`)

			if !re.MatchString(role) {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(
					e.Message {
						Message: "You are nor admin",
					},
				)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func GetRoleFromCTX(ctx context.Context) (string) {
			_role := ctx.Value(conextvalue.Role)
			role, ok := _role.(string)
			if !ok {
				log.WithFields(
					log.Fields{
						"package": "middleware/auth",
						"func": "GetRoleFromCTX",
						"err": "Failed to cast role to string",
					},
				).Panic()
			}

			return role
}

type getRoleFromClaim func(map[string]interface{}) (string, error)

func NewRoleGetter(
	claimName string,
	rolesSet map[string]struct{},
) getRoleFromClaim {
	return func(claims map[string]interface{}) (string, error) {
		claim, find := claims[claimName]

		if !find {
			return "", fmt.Errorf("Failed to get itlab claim")
		}

		_roles, ok := claim.([]interface{})
		if !ok {
			return "", fmt.Errorf("Failed to cast types")
		}

		roles := sliceOfInterfaceToSliceOfString(_roles)

		for _, role := range roles {
			if _, find := rolesSet[role]; find {
				return role, nil
			}
		}
		
		return "", fmt.Errorf("Failed to get rolse")
	}
}

func sliceOfInterfaceToSliceOfString(slice []interface{}) []string {
	var strs []string

	for _, elem := range slice {
		strs = append(strs, fmt.Sprint(elem))
	}

	return strs
}