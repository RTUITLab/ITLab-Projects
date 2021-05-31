package auth

import (
	"regexp"
	"github.com/ITLab-Projects/pkg/conextvalue/rolecontext"
	"context"
	"fmt"
	"net/http"

	ctxtoken "github.com/ITLab-Projects/pkg/conextvalue/token"
	"github.com/ITLab-Projects/pkg/statuscode"

	"github.com/auth0-community/go-auth0"
	"github.com/go-kit/kit/endpoint"
	"gopkg.in/square/go-jose.v2"

	"github.com/ITLab-Projects/pkg/config"
	log "github.com/sirupsen/logrus"
)

func NewGoKitAuth(
	cfg 	*config.AuthConfig,
	f		getRoleFromClaim,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("auth middleware")
			client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: cfg.KeyURL}, nil)
			configuration := auth0.NewConfiguration(client, []string{cfg.Audience}, cfg.Issuer, jose.RS256)
			validator := auth0.NewValidator(
				configuration, 
				nil,
			)

			_t, _ := ctxtoken.GetTokenFromContext(ctx)
			r := &http.Request{
				Header: http.Header{
					"Authorization": []string{_t},
				},
			}

			token, err := validator.ValidateRequest(
				r,
			)
			if err != nil {
				log.WithFields(log.Fields{
					"requiredAlgorithm" : "RS256",
					"error" : err,
				}).Debug("Token is not valid!")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Token is not valid"),
					http.StatusUnauthorized,
				)
			}
			claims := map[string]interface{}{}
			if err = validator.Claims(r, token, &claims); err != nil {
				log.WithFields(log.Fields{
					"requiredClaims" : "iss, aud, sub, role",
					"error" : err,
				}).Debug("Invalid claims!")
	
				
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Invalid claims"),
					http.StatusUnauthorized,
				)
			}

			role, err := f(claims)
			if err != nil {
				log.WithFields(log.Fields{
					"package" : "middleware/auth",
					"func": "authMiddleware",
					"error" : err,
				}).Debug("Failed to get role")

				return nil, statuscode.WrapStatusError(
					fmt.Errorf("Faield to get role"),
					http.StatusUnauthorized,
				)
			}

			ctx = rolecontext.New(
				ctx,
				role,
			)

			return next(ctx, request)
		}
	}
}

func EndpointAdminMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("admin middleware")
			role, err := rolecontext.GetRoleFromContext(ctx)
			if err != nil {
				log.WithFields(
					log.Fields{
						"package": "middleware/auth",
						"err": err,
					},
				).Panic()
			}

			re := regexp.MustCompile(`\w+.admin`)

			if !re.MatchString(role) {				
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("You are not admin"),
					http.StatusForbidden,
				)
			}

			return next(ctx, request)
		}
	}
}