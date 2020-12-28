package transport

import (
	"net/http"

	"github.com/StevenRojas/goaccess/pkg/entities"

	"github.com/StevenRojas/goaccess/pkg/codec"
	conf "github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/middlewares"

	"github.com/StevenRojas/goaccess/pkg/endpoints"
	"github.com/StevenRojas/goaccess/pkg/service"
	gokitJWT "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	gokitHTTP "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var corsMethods = []string{
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodGet,
	http.MethodDelete,
}

// Set API routes
// 	enpoints
// 	cors methods
// 	jwt middleware
// 	paths

func MakeHTTPHandlerForAccess(r *mux.Router, svc service.AccessService, config conf.SecurityConfig, logger conf.LoggerWrapper) {
	// create endpoints
	e := endpoints.MakeAccessEndpoints(
		svc,
		[]endpoint.Middleware{
			middlewares.JWTCheck(logger),
		},
	)
	// Apply CORS policy middleware
	r.Use(middlewares.CORSPolicies(corsMethods))
	r.Use(middlewares.ContentTypeMiddleware)
	// JWT decoder middleware
	//jwtDecoder, err := middlewares.DecodeJWT(jwt.SigningMethodHS256, config.JWTSecret, logger)
	// if err != nil {
	// 	logger.Error("invalid JWT", err)
	// }
	// Define server options to handle errors and decode JWT
	options := []gokitHTTP.ServerOption{
		gokitHTTP.ServerErrorEncoder(codec.HTTPErrorEncoder(logger)),
		gokitHTTP.ServerBefore(gokitJWT.HTTPToContext()),
		//gokitHTTP.ServerBefore(jwtDecoder),
	}
	// Initialize request validator
	entities.InitValidator()

	r.Methods(http.MethodPost).Path(getAccessPath("addRole")).Handler(gokitHTTP.NewServer(
		e.AddRole,
		codec.DecodeRegisterUserRequest,
		codec.JSONEncoder(logger),
		options...,
	))

}
