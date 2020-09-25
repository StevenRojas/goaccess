package service

import (
	"context"

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
	// AssignModule assign a module to a role
	AssignModule(ctx context.Context, name string, roleID string) error
	// UnassignModule unassign a module from a role
	UnassignModule(ctx context.Context, name string, roleID string) error
	// AssignSubModule assign a sub module to a role
	AssignSubModule(ctx context.Context, module string, name string, roleID string) error
	// UnassignSubModule unassign a sub module from a role
	UnassignSubModule(ctx context.Context, module string, name string, roleID string) error
	// AssignSection assign a section to a role
	AssignSection(ctx context.Context, module string, submodule string, name string, roleID string) error
	// UnassignSection unassign a section from a role
	UnassignSection(ctx context.Context, module string, submodule string, name string, roleID string) error
}

type access struct {
	repo repository.ModulesRepository
}

// NewAccessService return a new access service instance
func NewAccessService(modulesRepo repository.ModulesRepository) AccessService {
	return &access{
		repo: modulesRepo,
	}
}

// AddRole add a role and return its ID
func (a *access) AddRole(ctx context.Context, name string) (string, error) {
	return "", nil
}

//EditRole edit the role name
func (a *access) EditRole(ctx context.Context, ID string, name string) error {
	return nil
}

// DeleteRole removes a role and its relation with users
func (a *access) DeleteRole(ctx context.Context, ID string) error {
	return nil
}

// AssignModule assign a module to a role
func (a *access) AssignModule(ctx context.Context, name string, roleID string) error {
	return nil
}

// UnassignModule unassign a module from a role
func (a *access) UnassignModule(ctx context.Context, name string, roleID string) error {
	return nil
}

// AssignSubModule assign a sub module to a role
func (a *access) AssignSubModule(ctx context.Context, module string, name string, roleID string) error {
	return nil
}

// UnassignSubModule unassign a sub module from a role
func (a *access) UnassignSubModule(ctx context.Context, module string, name string, roleID string) error {
	return nil
}

// AssignSection assign a section to a role
func (a *access) AssignSection(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}

// UnassignSection unassign a section from a role
func (a *access) UnassignSection(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}
