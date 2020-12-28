package endpoints

import (
	"context"

	"github.com/StevenRojas/goaccess/pkg/codec"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AccessEndpoints struct {
	AddRole            endpoint.Endpoint
	EditRole           endpoint.Endpoint
	DeleteRole         endpoint.Endpoint
	AssignModules      endpoint.Endpoint
	UnassignModules    endpoint.Endpoint
	AssignSubModules   endpoint.Endpoint
	UnassignSubModules endpoint.Endpoint
	AssignSections     endpoint.Endpoint
	UnassignSections   endpoint.Endpoint
}

func MakeAccessEndpoints(
	s service.AccessService,
	middlewares []endpoint.Middleware) AccessEndpoints {
	return AccessEndpoints{
		AddRole:            wrapMiddlewares(makeAddRole(s), middlewares),
		EditRole:           wrapMiddlewares(makeEditRole(s), middlewares),
		DeleteRole:         wrapMiddlewares(makeDeleteRole(s), middlewares),
		AssignModules:      wrapMiddlewares(makeAssignModules(s), middlewares),
		UnassignModules:    wrapMiddlewares(makeUnassignModules(s), middlewares),
		AssignSubModules:   wrapMiddlewares(makeAssignSubModules(s), middlewares),
		UnassignSubModules: wrapMiddlewares(makeUnassignSubModules(s), middlewares),
		AssignSections:     wrapMiddlewares(makeAssignSections(s), middlewares),
		UnassignSections:   wrapMiddlewares(makeUnassignSections(s), middlewares),
	}
}

func wrapMiddlewares(ep endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for i := range middlewares {
		ep = middlewares[i](ep)
	}
	return ep
}

func makeAddRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// user, ok := request.(entities.User)
		// if !ok {
		// 	return nil, e.HTTPBadRequest(errors.New("unable to cast the request to UserRequest"))
		// }
		//err := s.Register(ctx, &user)
		return &codec.EmptyResponse{}, nil
	}
}

func makeEditRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeDeleteRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeAssignModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeAssignSubModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignSubModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeAssignSections(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignSections(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return &codec.EmptyResponse{}, nil
	}
}
