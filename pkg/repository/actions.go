package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/go-redis/redis/v8"
)

// ActionsRepository modules repository
type ActionsRepository interface {
	// AssignActions assign actions to a role
	AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// UnassignActions unassign actions from a role
	UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error
	// GetActionListByModule get a json list with the actions can be performed by a user in a module
	GetActionListByModule(ctx context.Context, module string, userID string) (string, error)
	// CheckPermission checks if a user has permission to perform an action
	CheckPermission(ctx context.Context, action string, userID string) (bool, error)
	// ActionsStructure returns a list of available actions grouped by module
	ActionsStructure(ctx context.Context) ([]entities.ActionModule, error)
}

type actionsRepo struct {
	c *redis.Client
}

// NewActionsRepository creates a new repository instance
func NewActionsRepository(ctx context.Context, client *redis.Client) (ActionsRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &actionsRepo{
		c: client,
	}, nil
}

// AssignActions assign actions to a role
func (r *actionsRepo) AssignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	key := fmt.Sprintf(roleActionsKey, rolesKey, roleID, module[:2], submodule)
	_, err := r.c.SAdd(ctx, key, actions).Result()
	return err
}

// UnassignActions unassign actions from a role
func (r *actionsRepo) UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	key := fmt.Sprintf(roleActionsKey, rolesKey, roleID, module[:2], submodule)
	_, err := r.c.SRem(ctx, key, actions).Result()
	return err
}

// GetActionListByModule get a json list with the actions can be performed by a user in a module
func (r *actionsRepo) GetActionListByModule(ctx context.Context, module string, userID string) (string, error) {
	key := fmt.Sprintf(actionsByModuleKey, userID, module)
	actions, err := r.c.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return actions, nil
}

// CheckPermission checks if a user has permission to perform an action
func (r *actionsRepo) CheckPermission(ctx context.Context, action string, userID string) (bool, error) {
	key := fmt.Sprintf(hasPesmissionKey, action)
	isMember, err := r.c.SIsMember(ctx, key, userID).Result()
	if err != nil {
		return false, err
	}
	return isMember, nil
}

// ActionsStructure returns a list of available modules, submodules and sections
func (r *actionsRepo) ActionsStructure(ctx context.Context) ([]entities.ActionModule, error) {
	keys, err := r.c.Keys(ctx, actionTemplateKey+"*").Result()
	if err != nil {
		return nil, err
	}
	var modules []entities.ActionModule
	var module entities.ActionModule
	for _, key := range keys {
		j, err := r.c.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(j), &module)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}
	return modules, err
}
