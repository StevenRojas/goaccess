package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// ModulesRepository modules repository
type ModulesRepository interface {
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

type modulesRepo struct {
	c *redis.Client
}

// NewModulesRepository creates a new repository instance
func NewModulesRepository(ctx context.Context, client *redis.Client) (ModulesRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &modulesRepo{
		c: client,
	}, nil
}

// AssignModule assign a module to a role
func (r *modulesRepo) AssignModule(ctx context.Context, name string, roleID string) error {
	return nil
}

// UnassignModule unassign a module from a role
func (r *modulesRepo) UnassignModule(ctx context.Context, name string, roleID string) error {
	return nil
}

// AssignSubModule assign a sub module to a role
func (r *modulesRepo) AssignSubModule(ctx context.Context, module string, name string, roleID string) error {
	return nil
}

// UnassignSubModule unassign a sub module from a role
func (r *modulesRepo) UnassignSubModule(ctx context.Context, module string, name string, roleID string) error {
	return nil
}

// AssignSection assign a section to a role
func (r *modulesRepo) AssignSection(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}

// UnassignSection unassign a section from a role
func (r *modulesRepo) UnassignSection(ctx context.Context, module string, submodule string, name string, roleID string) error {
	return nil
}
