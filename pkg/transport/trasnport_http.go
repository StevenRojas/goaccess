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

	r.Methods(http.MethodGet).Path(getAccessPath("listRoles")).Handler(gokitHTTP.NewServer(
		e.ListRoles,
		codec.DecodePaginatedListRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getAccessPath("rolesByUser")).Handler(gokitHTTP.NewServer(
		e.ListRolesByUser,
		codec.DecodeGetRolesByUserRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getAccessPath("addRole")).Handler(gokitHTTP.NewServer(
		e.AddRole,
		codec.DecodeAddRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPut).Path(getAccessPath("editRole")).Handler(gokitHTTP.NewServer(
		e.EditRole,
		codec.DecodeEditRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getAccessPath("deleteRole")).Handler(gokitHTTP.NewServer(
		e.DeleteRole,
		codec.DecodeDeleteRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getAccessPath("getAllModules")).Handler(gokitHTTP.NewServer(
		e.GetAllModules,
		codec.DecodeEmptyRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getAccessPath("getAssignedModules")).Handler(gokitHTTP.NewServer(
		e.GetAssignedModules,
		codec.DecodeGetRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getAccessPath("assignModules")).Handler(gokitHTTP.NewServer(
		e.AssignModules,
		codec.DecodeAssignModulesRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getAccessPath("unassignModules")).Handler(gokitHTTP.NewServer(
		e.UnassignModules,
		codec.DecodeUnassignModulesRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getAccessPath("getAssignedSubModules")).Handler(gokitHTTP.NewServer(
		e.GetAssignedSubModules,
		codec.DecodeGetRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getAccessPath("assignSubModules")).Handler(gokitHTTP.NewServer(
		e.AssignSubModules,
		codec.DecodeAssignSubModulesRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getAccessPath("unassignSubModules")).Handler(gokitHTTP.NewServer(
		e.UnassignSubModules,
		codec.DecodeUnassignSubModulesRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getAccessPath("getAssignedSections")).Handler(gokitHTTP.NewServer(
		e.GetAssignedSections,
		codec.DecodeGetRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getAccessPath("assignSections")).Handler(gokitHTTP.NewServer(
		e.AssignSections,
		codec.DecodeAssignSectionsRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getAccessPath("unassignSections")).Handler(gokitHTTP.NewServer(
		e.UnassignSections,
		codec.DecodeUnassignSectionsRequest,
		codec.JSONEncoder(logger),
		options...,
	))
}

func MakeHTTPHandlerForActions(r *mux.Router, svc service.AuthorizationService, config conf.SecurityConfig, logger conf.LoggerWrapper) {
	// create endpoints
	e := endpoints.MakeActionsEndpoints(
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
	// entities.InitValidator()

	r.Methods(http.MethodGet).Path(getActionsPath("listUsers")).Handler(gokitHTTP.NewServer(
		e.ListUsers,
		codec.DecodePaginatedListRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getActionsPath("usersByRole")).Handler(gokitHTTP.NewServer(
		e.ListUsersByRole,
		codec.DecodeGetRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getActionsPath("assignRole")).Handler(gokitHTTP.NewServer(
		e.AssignRole,
		codec.DecodeAssignUnassignRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getActionsPath("unassignRole")).Handler(gokitHTTP.NewServer(
		e.UnassignRole,
		codec.DecodeAssignUnassignRoleRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodPost).Path(getActionsPath("assignActions")).Handler(gokitHTTP.NewServer(
		e.AssignActions,
		codec.DecodeAssignActionsRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodDelete).Path(getActionsPath("unassignActions")).Handler(gokitHTTP.NewServer(
		e.UnassignActions,
		codec.DecodeUnassignActionsRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getActionsPath("getAccessList")).Handler(gokitHTTP.NewServer(
		e.GetAccessList,
		codec.DecodeGetAccessRequest,
		codec.JSONEncoder(logger),
		options...,
	))

	r.Methods(http.MethodGet).Path(getActionsPath("getActionList")).Handler(gokitHTTP.NewServer(
		e.GetActionList,
		codec.DecodeGetActionsRequest,
		codec.JSONEncoder(logger),
		options...,
	))

}
