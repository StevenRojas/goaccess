package main

import (
	"context"
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/service"
)

func main() {
	ctx := context.Background()
	serviceConfig, err := configuration.Read()
	if err != nil {
		panic(err)
	}
	logger := configuration.NewLogger(serviceConfig.Server)
	logger.Debug("creating services...")
	factory := service.NewServiceFactory(ctx, serviceConfig)
	factory.Setup()
	authenticationService := factory.CreateAuthenticationService()
	accessService := factory.CreateAccessService()
	authorizationService := factory.CreateAuthorizationService()
	initService := factory.CreateInitializationService()
	logger.Debug("services ready")
	fmt.Printf("%T --- %T ---%T ---%T ---\n\n", authenticationService, accessService, authorizationService, initService)

	initService.Init(ctx, false)

}
