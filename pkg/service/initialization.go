package service

import (
	"context"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"
	"github.com/StevenRojas/goaccess/pkg/utils"
)

// InitializationService initialization service
type InitializationService interface {
	// Init initialize modules, sections and actions if needed
	Init(ctx context.Context, force bool) error
}

type initialization struct {
	repo repository.InitRepository
	jh   utils.JSONHandler
}

// NewInitService return a new authorization service instance
func NewInitService(initRepo repository.InitRepository, jsonHandler utils.JSONHandler) InitializationService {
	return &initialization{
		repo: initRepo,
		jh:   jsonHandler,
	}
}

// Init initialize modules, sections and actions if needed
func (i *initialization) Init(ctx context.Context, force bool) error {
	if !force {
		isset, err := i.repo.IsSetConfig(ctx)
		if err != nil {
			return err
		}
		if isset {
			return nil
		}
	}
	i.repo.UnsetConfig(ctx)
	modules, err := i.jh.Modules()
	if err != nil {
		return err
	}

	err = i.initModules(ctx, modules)
	if err != nil {
		return err
	}

	err = i.repo.SetAsConfigured(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *initialization) initModules(ctx context.Context, modules []entities.ModuleInit) error {
	var names []string
	for _, module := range modules {
		names = append(names, module.Name)
		i.repo.AddModule(ctx, module)
	}
	return nil
}
