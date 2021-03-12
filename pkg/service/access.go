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
	// IsRoleExist check if the role exists
	IsRoleExist(ctx context.Context, ID string) (bool, error)
	//EditRole edit the role name
	EditRole(ctx context.Context, ID string, name string) error
	// DeleteRole removes a role and its relation with users
	DeleteRole(ctx context.Context, ID string) error
	// // GetAllModules get a list of available modules
	// GetAllModules(ctx context.Context) error
	// // GetAssignedModules get assign modules to a role
	// GetAssignedModules(ctx context.Context, roleID string) error
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
	// ModulesListByRole returns a list of available modules for a given role
	ModulesListByRole(ctx context.Context, roleID string) ([]string, error)
	// SubModulesListByRole returns a list of available submodules for a given role
	SubModulesListByRole(ctx context.Context, roleID string) (map[string][]string, error)
	// SectionsListByRole returns a list of available sections for a given role
	SectionsListByRole(ctx context.Context, roleID string) (map[string]map[string][]string, error)
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

// IsRoleExist check if the role exists
func (a *access) IsRoleExist(ctx context.Context, ID string) (bool, error) {
	return a.rolesRepo.IsValidRole(ctx, ID)
}

// DeleteRole removes a role and its relation with users
func (a *access) DeleteRole(ctx context.Context, ID string) error {
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

// ModulesListByRole returns a list of available modules for a given role
func (a *access) ModulesListByRole(ctx context.Context, roleID string) ([]string, error) {
	return a.modulesRepo.ModulesListByRole(ctx, roleID)
}

// SubModulesListByRole returns a list of available submodules for a given role
func (a *access) SubModulesListByRole(ctx context.Context, roleID string) (map[string][]string, error) {
	return a.modulesRepo.SubModulesListByRole(ctx, roleID)
}

// SectionsListByRole returns a list of available sections for a given role
func (a *access) SectionsListByRole(ctx context.Context, roleID string) (map[string]map[string][]string, error) {
	return a.modulesRepo.SectionsListByRole(ctx, roleID)
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
