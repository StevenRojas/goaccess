package endpoints

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess/pkg/codec"
	"github.com/StevenRojas/goaccess/pkg/entities"
	e "github.com/StevenRojas/goaccess/pkg/errors"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type ActionsEndpoints struct {
	ListUsers       endpoint.Endpoint
	ListUsersByRole endpoint.Endpoint
	GetAccessList   endpoint.Endpoint
	GetActionList   endpoint.Endpoint
	AssignActions   endpoint.Endpoint
	UnassignActions endpoint.Endpoint
	AssignRole      endpoint.Endpoint
	UnassignRole    endpoint.Endpoint
}

func MakeActionsEndpoints(
	s service.AuthorizationService,
	middlewares []endpoint.Middleware) ActionsEndpoints {
	return ActionsEndpoints{
		ListUsers:       wrapMiddlewares(makeListUsers(s), middlewares),
		ListUsersByRole: wrapMiddlewares(makeListUsersByRole(s), middlewares),
		GetAccessList:   wrapMiddlewares(makeGetAccessList(s), middlewares),
		GetActionList:   wrapMiddlewares(makeGetActionList(s), middlewares),
		AssignActions:   wrapMiddlewares(makeAssignActions(s), middlewares),
		UnassignActions: wrapMiddlewares(makeUnassignActions(s), middlewares),
		AssignRole:      wrapMiddlewares(makeAssignRole(s), middlewares),
		UnassignRole:    wrapMiddlewares(makeUnassignRole(s), middlewares),
	}
}

func makeListUsers(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.ListUsers(ctx)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get a list of users", err)
		}
		return users, nil
	}
}

func makeListUsersByRole(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		users, err := s.ListUsersByRole(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get a list of users", err)
		}
		return users, nil
	}
}

func makeGetAccessList(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		accessList, err := s.GetAccessList(ctx, userID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get access list", err)
		}
		return accessList, nil
	}
}

func makeGetActionList(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(map[string]string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to map"))
		}
		actionList, err := s.GetActionListByModule(ctx, req["module"], req["userID"])
		if err != nil {
			return nil, e.HTTPConflict("Unable to get access list", err)
		}
		return actionList, nil
	}
}

func makeAssignActions(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		actions, ok := request.(entities.ActionList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to ActionList"))
		}
		err := s.AssignActions(ctx, actions.RoleID, actions.Module, actions.SubModule, actions.Actions)
		if err != nil {
			return nil, e.HTTPConflict("Unable to add role", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignActions(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		actions, ok := request.(entities.ActionList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to ActionList"))
		}
		err := s.UnassignActions(ctx, actions.RoleID, actions.Module, actions.SubModule, actions.Actions)
		if err != nil {
			return nil, e.HTTPConflict("Unable to add role", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeAssignRole(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		params, ok := request.(map[string]string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request"))
		}
		roleID := params["roleID"]
		userID := params["userID"]
		err := s.AssignRole(ctx, userID, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to assign role", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignRole(s service.AuthorizationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		params, ok := request.(map[string]string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request"))
		}
		roleID := params["roleID"]
		userID := params["userID"]
		err := s.UnassignRole(ctx, userID, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to unassign role", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}
