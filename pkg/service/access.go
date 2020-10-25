package service

import (
	"context"
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/events"
	"github.com/StevenRojas/goaccess/pkg/repository"
)

// AccessService access service to handle modules, submodules and sections
type AccessService interface {
	// AddRole add a role and return its ID
	AddRole(ctx context.Context, name string) (string, error)
	//EditRole edit the role name
	EditRole(ctx context.Context, ID string, name string) error
	// DeleteRole removes a role and its relation with users
	DeleteRole(ctx context.Context, ID string) error
	// AssignModules assign modules to a role
	AssignModules(ctx context.Context, roleID string, modules []string) error
	// UnassignModules unassign modules from a role
	UnassignModules(ctx context.Context, roleID string, modules []string) error
	// AssignSubModules assign submodules to a role
	AssignSubModules(ctx context.Context, roleID string, module string, submodules []string) error
	// UnassignSubModules unassign submodules from a role
	UnassignSubModules(ctx context.Context, roleID string, module string, submodules []string) error
	// AssignSections assign sections to a role
	AssignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error
	// UnassignSections unassign sections from a role
	UnassignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error
	// ModulesForNewRole returns a list of available modules for create a new role
	ModulesForNewRole(ctx context.Context) ([]entities.Module, error)
	// ActionsForNewRole returns a list of available actions for create a new role
	ActionsForNewRole(ctx context.Context) ([]entities.ActionModule, error)
	// GetRoleAccessList get a json of modules, submodules and sections for the given role
	GetRoleAccessList(ctx context.Context, roleID string) (string, error)
}

type access struct {
	modulesRepo    repository.ModulesRepository
	rolesRepo      repository.RolesRepository
	actionsRepo    repository.ActionsRepository
	subscriberFeed events.SubscriberFeed
}

// NewAccessService return a new access service instance
func NewAccessService(
	modulesRepo repository.ModulesRepository,
	rolesRepo repository.RolesRepository,
	actionsRepo repository.ActionsRepository,
	subscriberFeed events.SubscriberFeed,
) AccessService {
	return &access{
		modulesRepo:    modulesRepo,
		rolesRepo:      rolesRepo,
		actionsRepo:    actionsRepo,
		subscriberFeed: subscriberFeed,
	}
}

// AddRole add a role and return its ID
func (a *access) AddRole(ctx context.Context, name string) (string, error) {
	return a.rolesRepo.AddRole(ctx, name)
}

//EditRole edit the role name
func (a *access) EditRole(ctx context.Context, ID string, name string) error {
	return a.rolesRepo.EditRole(ctx, ID, name)
}

// DeleteRole removes a role and its relation with users
func (a *access) DeleteRole(ctx context.Context, ID string) error {
	// TODO: Update access for assigned users and role access/actions definitions
	// roleEvent := &entities.RoleEvent{RoleID: 13}
	// go a.subscriberFeed.Send(roleEvent)
	return a.rolesRepo.DeleteRole(ctx, ID)
}

// AssignModules assign a module to a role
func (a *access) AssignModules(ctx context.Context, roleID string, modules []string) error {
	for _, module := range modules {
		err := a.modulesRepo.AssignModule(ctx, roleID, module)
		if err != nil {
			return err
		}
	}
	// TODO: Update access for assigned users and role access/actions definitions
	// roleEvent := &entities.RoleEvent{RoleID: 13}
	// go a.subscriberFeed.Send(roleEvent)
	return nil
}

// UnassignModules unassign a module from a role
func (a *access) UnassignModules(ctx context.Context, roleID string, modules []string) error {
	for _, module := range modules {
		err := a.modulesRepo.UnassignModule(ctx, roleID, module)
		if err != nil {
			return err
		}
	}
	return nil
}

// AssignSubModules assign a sub module to a role
func (a *access) AssignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	err := a.modulesRepo.AssignSubModules(ctx, roleID, module, submodules)
	if err != nil {
		return err
	}
	// Update access for assigned user
	return nil
}

// UnassignSubModules unassign a sub module from a role
func (a *access) UnassignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	err := a.modulesRepo.UnassignSubModules(ctx, roleID, module, submodules)
	if err != nil {
		return err
	}
	// Update access for assigned user
	return nil
}

// AssignSections assign a section to a role
func (a *access) AssignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	err := a.modulesRepo.AssignSections(ctx, roleID, module, submodule, sections)
	if err != nil {
		return err
	}
	// Update access for assigned user
	return nil
}

// UnassignSections unassign a section from a role
func (a *access) UnassignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	err := a.modulesRepo.UnassignSections(ctx, roleID, module, submodule, sections)
	if err != nil {
		return err
	}
	// Update access for assigned user
	return nil
}

// ModulesForNewRole returns a list of available modules for create a new role
func (a *access) ModulesForNewRole(ctx context.Context) ([]entities.Module, error) {
	return a.modulesRepo.ModulesStructure(ctx)
}

// ActionsForNewRole returns a list of available actions for create a new role
func (a *access) ActionsForNewRole(ctx context.Context) ([]entities.ActionModule, error) {
	return a.actionsRepo.ActionsStructure(ctx)
}

// GetRoleAccessList get a json of modules, submodules and sections for the given role
func (a *access) GetRoleAccessList(ctx context.Context, roleID string) (string, error) {
	// Get modules, submodules and sections assigned to the role
	assignations, err := a.modulesRepo.AssignationsByRole(ctx, roleID)
	if err != nil {
		return "", err
	}
	// Get the module structure
	moduleStructure, err := a.modulesRepo.ModuleStructure(ctx, roleID)
	if err != nil {
		return "", err
	}
	// Set the corresponding access to the structure

	fmt.Printf("%v\n%v\n", assignations, moduleStructure)
	return "", nil
}
