package service

import (
	"context"

	"github.com/StevenRojas/goaccess/pkg/repository"
)

// AuthorizationService authorization service to handle modules, submodules and sections
type AuthorizationService interface {
	// AssignAction assign an action to a role
	AssignAction(ctx context.Context, module string, submodule string, name string, roleID string) error
	// UnassignAction unassign an action from a role
	UnassignAction(ctx context.Context, module string, submodule string, name string, roleID string) error
	// AssingRole assign a role to a user
	AssignRole(ctx context.Context, userID string, roleID string) error
	// UnassignRole unassign a role from a user
	UnassignRole(ctx context.Context, userID string, roleID string) error

	// GetAccessList get a json of modules, submodules and sections where the user has access
	GetAccessList(ctx context.Context, userID string) (string, error)
	// GetActionListByModule get a json list with the actions can be performed by a user in a module
	GetActionListByModule(ctx context.Context, module string, userID string) (string, error)
	// CheckPermission checks if a user has permission to perform an action
	CheckPermission(ctx context.Context, action string, userID string) (bool, error)
}

type authorization struct {
	repo repository.RolesRepository
}

// NewAuthorizationService return a new authorization service instance
func NewAuthorizationService(rolesRepo repository.RolesRepository) AuthorizationService {
	return &authorization{
		repo: rolesRepo,
	}
}

// AssignAction assign an action to a role
func (a *authorization) AssignAction(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}

// UnassignAction unassign an action from a role
func (a *authorization) UnassignAction(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}

// AssingRole assign a role to a user
func (a *authorization) AssignRole(ctx context.Context, userID string, roleID string) error {
	return nil
}

// UnassignRole unassign a role from a user
func (a *authorization) UnassignRole(ctx context.Context, userID string, roleID string) error {
	return nil
}

// GetAccessList get a json of modules, submodules and sections where the user has access
func (a *authorization) GetAccessList(ctx context.Context, userID string) (string, error) {
	return "", nil
}

// GetActionListByModule get a json list with the actions can be performed by a user in a module
func (a *authorization) GetActionListByModule(ctx context.Context, module string, userID string) (string, error) {
	return "", nil
}

// CheckPermission checks if a user has permission to perform an action
func (a *authorization) CheckPermission(ctx context.Context, action string, userID string) (bool, error) {
	return false, nil
}
