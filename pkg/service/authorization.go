package service

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess/pkg/repository"
)

// AuthorizationService authorization service to handle modules, submodules and sections
type AuthorizationService interface {
	// AssignActions assign actions to a role
	AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// UnassignActions unassign actions from a role
	UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// AssingRoles assign roles to a user
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	// UnassignRoles unassign roles from a user
	UnassignRoles(ctx context.Context, userID string, roleIDs []string) error

	// GetRoleActionList get a json list with the actions for the given role
	GetRoleActionList(ctx context.Context, module string, roleID string) (string, error)
	// GetAccessList get a json of modules, submodules and sections where the user has access
	GetAccessList(ctx context.Context, userID string) (string, error)
	// GetActionListByModule get a json list with the actions can be performed by a user in a module
	GetActionListByModule(ctx context.Context, module string, userID string) (string, error)
	// CheckPermission checks if a user has permission to perform an action
	CheckPermission(ctx context.Context, action string, userID string) (bool, error)
}

type authorization struct {
	modulesRepo repository.ModulesRepository
	rolesRepo   repository.RolesRepository
	actionsRepo repository.ActionsRepository
}

// NewAuthorizationService return a new authorization service instance
func NewAuthorizationService(
	modulesRepo repository.ModulesRepository,
	rolesRepo repository.RolesRepository,
	actionsRepo repository.ActionsRepository) AuthorizationService {
	return &authorization{
		modulesRepo: modulesRepo,
		rolesRepo:   rolesRepo,
		actionsRepo: actionsRepo,
	}
}

// AssignActions assign actions to a role
func (a *authorization) AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	err := a.actionsRepo.AssignActions(ctx, roleID, module, submodule, actions)
	if err != nil {
		return err
	}
	// Update acctions for assigned user
	return nil
}

// UnassignActions unassign actions from a role
func (a *authorization) UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	err := a.actionsRepo.UnassignActions(ctx, roleID, module, submodule, actions)
	if err != nil {
		return err
	}
	// Update acctions for assigned user
	return nil
}

// AssingRoles assign roles to a user
func (a *authorization) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	err := a.rolesRepo.AssignRoles(ctx, userID, roleIDs)
	if err != nil {
		return err
	}
	// Call to roles updater
	return nil
}

// UnassignRoles unassign roles from a user
func (a *authorization) UnassignRoles(ctx context.Context, userID string, roleIDs []string) error {
	err := a.rolesRepo.UnassignRoles(ctx, userID, roleIDs)
	if err != nil {
		return err
	}
	// Call to roles updater
	return nil
}

// GetRoleActionList get a json list with the actions for the given role
func (a *authorization) GetRoleActionList(ctx context.Context, module string, roleID string) (string, error) {
	return "", nil
}

// GetAccessList get a json of modules, submodules and sections where the user has access
func (a *authorization) GetAccessList(ctx context.Context, userID string) (string, error) {
	access, err := a.modulesRepo.GetAccessList(ctx, userID)
	if err != nil {
		return "", err
	}
	if access == "" {
		return "", errors.New("User has not access defined")
	}
	return access, nil
}

// GetActionListByModule get a json list with the actions can be performed by a user in a module
func (a *authorization) GetActionListByModule(ctx context.Context, module string, userID string) (string, error) {
	actions, err := a.actionsRepo.GetActionListByModule(ctx, module, userID)
	if err != nil {
		return "", err
	}
	if actions == "" {
		return "", errors.New("User has not actions defined")
	}
	return actions, nil
}

// CheckPermission checks if a user has permission to perform an action
func (a *authorization) CheckPermission(ctx context.Context, action string, userID string) (bool, error) {
	return a.actionsRepo.CheckPermission(ctx, action, userID)
}
