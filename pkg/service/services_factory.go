package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/events"
	"github.com/StevenRojas/goaccess/pkg/repository"
	"github.com/StevenRojas/goaccess/pkg/utils"
	"github.com/go-redis/redis/v8"
)

// ServicesFactory service factory interface
type ServicesFactory interface {
	// Setup repositories and event listeners
	Setup()
	// CreateAuthenticationService create Authentication service
	CreateAuthenticationService() AuthenticationService
	// CreateAccessService create Access service
	CreateAccessService() AccessService
	// CreateAuthorizationService create Authorization service
	CreateAuthorizationService() AuthorizationService
	// CreateInitService create Initialization service
	CreateInitializationService() InitializationService
}

type serviceFactory struct {
	reposReady     bool
	ctx            context.Context
	serviceConfig  *configuration.ServiceConfig
	usersRepo      repository.UsersRepository
	modulesRepo    repository.ModulesRepository
	rolesRepo      repository.RolesRepository
	actionsRepo    repository.ActionsRepository
	initRepo       repository.InitRepository
	subscriberFeed events.SubscriberFeed
}

// NewServiceFactory get a new service factory instance
func NewServiceFactory(ctx context.Context, serviceConfig *configuration.ServiceConfig) ServicesFactory {
	return &serviceFactory{
		reposReady:    false,
		ctx:           ctx,
		serviceConfig: serviceConfig,
	}
}

// Setup repositories and event listeners
func (sb *serviceFactory) Setup() {
	fmt.Println(sb.serviceConfig.Redis.DB)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     sb.serviceConfig.Redis.Addr,
		Password: sb.serviceConfig.Redis.Pass,
		DB:       sb.serviceConfig.Redis.DB,
	})

	var err error
	sb.usersRepo, err = repository.NewUsersRepository(sb.ctx, redisClient)
	if err != nil {
		panic(errors.New("Unable to create users repository"))
	}
	sb.modulesRepo, err = repository.NewModulesRepository(sb.ctx, redisClient)
	if err != nil {
		panic(errors.New("Unable to create modules repository"))
	}
	sb.rolesRepo, err = repository.NewRolesRepository(sb.ctx, redisClient)
	if err != nil {
		panic(errors.New("Unable to create roles repository"))
	}
	sb.actionsRepo, err = repository.NewActionsRepository(sb.ctx, redisClient)
	if err != nil {
		panic(errors.New("Unable to create actions repository"))
	}
	sb.initRepo, err = repository.NewInitRepository(sb.ctx, redisClient)
	if err != nil {
		panic(errors.New("Unable to create init repository"))
	}
	sb.subscriberFeed = events.NewSubscriber()
	sb.reposReady = true
}

// CreateAuthenticationService create Authentication service
func (sb serviceFactory) CreateAuthenticationService() AuthenticationService {
	if !sb.reposReady {
		panic(errors.New("Repositories not created, use Setup method first"))
	}
	jwtHander := utils.NewJwtHandler(sb.serviceConfig.Security)
	return NewAuthenticationService(sb.usersRepo, jwtHander)
}

// CreateAccessService create Access service
func (sb serviceFactory) CreateAccessService() AccessService {
	if !sb.reposReady {
		panic(errors.New("Repositories not created, use Setup method first"))
	}
	// Role events listener
	accessListener := events.NewAccessListener(sb.modulesRepo, sb.rolesRepo, sb.subscriberFeed)
	go accessListener.RegisterAccessListener()
	actionListener := events.NewActionListener(sb.actionsRepo, sb.rolesRepo, sb.subscriberFeed)
	go actionListener.RegisterActionListener()

	return NewAccessService(sb.modulesRepo, sb.rolesRepo, sb.actionsRepo, sb.subscriberFeed)
}

// CreateAuthorizationService create Authorization service
func (sb serviceFactory) CreateAuthorizationService() AuthorizationService {
	if !sb.reposReady {
		panic(errors.New("Repositories not created, use Setup method first"))
	}
	return NewAuthorizationService(sb.modulesRepo, sb.rolesRepo, sb.actionsRepo, sb.usersRepo, sb.subscriberFeed)
}

// CreateInitService create Initialization service
func (sb serviceFactory) CreateInitializationService() InitializationService {
	if !sb.reposReady {
		panic(errors.New("Repositories not created, use Setup method first"))
	}
	path, _ := os.Getwd()
	path = path + "/init"
	jsonHandler := utils.NewJSONHandler(path)
	return NewInitService(sb.initRepo, jsonHandler)
}
