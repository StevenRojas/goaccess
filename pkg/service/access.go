package service

import (
	"context"
	"encoding/json"
	"errors"

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
	// ModulesList returns a list of available modules
	ModulesList(ctx context.Context) ([]string, error)
	// ModuleStructure returns the module structure to create a new role
	ModuleStructure(ctx context.Context, name string) (*entities.Module, error)
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
	if ok, _ := a.rolesRepo.IsValidRole(ctx, ID); !ok {
		return errors.New("Role not found")
	}
	err := a.rolesRepo.DeleteRole(ctx, ID)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: ID}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// AssignModules assign a module to a role
func (a *access) AssignModules(ctx context.Context, roleID string, modules []string) error {
	if ok, _ := a.rolesRepo.IsValidRole(ctx, roleID); !ok {
		return errors.New("Role not found")
	}
	for _, module := range modules {
		err := a.modulesRepo.AssignModule(ctx, roleID, module)
		if err != nil {
			return err
		}
	}
	// Update access for assigned users
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
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
	// Update access for assigned users
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// AssignSubModules assign a sub module to a role
func (a *access) AssignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	err := a.modulesRepo.AssignSubModules(ctx, roleID, module, submodules)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// UnassignSubModules unassign a sub module from a role
func (a *access) UnassignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	err := a.modulesRepo.UnassignSubModules(ctx, roleID, module, submodules)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// AssignSections assign a section to a role
func (a *access) AssignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	err := a.modulesRepo.AssignSections(ctx, roleID, module, submodule, sections)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// UnassignSections unassign a section from a role
func (a *access) UnassignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	err := a.modulesRepo.UnassignSections(ctx, roleID, module, submodule, sections)
	if err != nil {
		return err
	}
	roleEvent := &entities.RoleEvent{RoleID: roleID, EventType: entities.EventTypeAccess}
	go a.subscriberFeed.Send(roleEvent)
	return nil
}

// ModulesList returns a list of available modules
func (a *access) ModulesList(ctx context.Context) ([]string, error) {
	return a.modulesRepo.ModulesList(ctx)
}

// ModuleStructure returns the module structure to create a new role
func (a *access) ModuleStructure(ctx context.Context, name string) (*entities.Module, error) {
	return a.modulesRepo.ModuleStructure(ctx, name)
}

// GetRoleAccessList get a json of modules, submodules and sections for the given role
func (a *access) GetRoleAccessList(ctx context.Context, roleID string) (string, error) {
	// Get modules, submodules and sections assigned to the role
	assignations, err := a.modulesRepo.AssignationsByRole(ctx, roleID)
	if err != nil {
		return "", err
	}
	j, err := json.Marshal(assignations)
	return string(j), err
}
