package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/StevenRojas/goaccess/pkg/transport"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
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

	router := mux.NewRouter()
	transport.MakeHTTPHandlerForAccess(router, accessService, serviceConfig.Security, logger)

	var runGroup group.Group
	{
		httpServer := http.Server{
			Addr:    serviceConfig.Server.HTTP,
			Handler: router,
		}
		runGroup.Add(func() error {
			logger.Info("HTTP server listen at " + serviceConfig.Server.HTTP)
			return httpServer.ListenAndServe() // TODO: support TLS
		}, func(err error) {
			httpServer.Shutdown(context.Background())
			logger.Error("HTTP server shutdown with error", "error", err)
		})
	}

	{
		cancel := make(chan struct{})
		runGroup.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case s := <-c:
				return fmt.Errorf("signal received %s", s)
			case <-cancel:
				return nil
			}
		}, func(error) {
			close(cancel)
		})
	}
	runGroup.Run()
	logger.Info("server terminated")
}
