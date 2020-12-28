package middlewares

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	gokitHTTP "github.com/go-kit/kit/transport/http"

	log "github.com/StevenRojas/goaccess/pkg/configuration"
)

// JWTCheck check if JWT is valid
func JWTCheck(logger log.LoggerWrapper) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger.Info("Validating JWT")
			return next(ctx, request)
		}
	}
}

// DecodeJWT decode JWT into claims
func DecodeJWT(signingMethod jwt.SigningMethod, secret string, logger log.LoggerWrapper) (gokitHTTP.RequestFunc, error) {
	logger.Info("decoding JWT")
	return nil, nil
}
