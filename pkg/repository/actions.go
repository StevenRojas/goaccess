package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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
	// SetActionList sets the action list for a given user based on all assigned roles
	SetActionList(ctx context.Context, userID string) error
	// ActionsByRole get a list of actions assigned to the role
	ActionsByRole(ctx context.Context, roleID string) (map[string]interface{}, error)
	// RemoveActionsByUser Remove action list for a given user
	RemoveActionsByUser(ctx context.Context, userID string) error
	// UpdateActionList update the list of actions to quick access while checking permissions
	UpdateActionList(ctx context.Context, roleID string) error
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
	key := fmt.Sprintf(roleActionsKey, roleID, module, submodule)
	_, err := r.c.SAdd(ctx, key, actions).Result()
	if err != nil {
		return err
	}
	return err
}

// UnassignActions unassign actions from a role
func (r *actionsRepo) UnassignActions(ctx context.Context, roleID string, module string, submodule string, actions []string) error {
	key := fmt.Sprintf(roleActionsKey, roleID, module, submodule)
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
	key := fmt.Sprintf(hasPesmissionKey, userID)
	isMember, err := r.c.SIsMember(ctx, key, action).Result()
	if err != nil {
		return false, err
	}
	return isMember, nil
}

// SetActionList sets the action list for a given user based on all assigned roles
func (r *actionsRepo) SetActionList(ctx context.Context, userID string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	roles, err := r.c.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	sort.Strings(roles)
	//assignations := make(map[string]interface{})
	for _, role := range roles {
		assignedModules, err := r.ActionsByRole(ctx, role)
		if err != nil {
			return err
		}
		// Modules are override with the latest role
		for module, assignation := range assignedModules {
			//assignations[module] = assignation
			j, _ := json.Marshal(assignation)
			key := fmt.Sprintf(actionsByModuleKey, userID, module)
			_, err = r.c.Set(ctx, key, j, 0).Result()
		}
	}
	if err != nil && err != redis.Nil {
		return err
	}
	return err
}

// ActionsByRole get a list of actions assigned to the role
func (r *actionsRepo) ActionsByRole(ctx context.Context, roleID string) (map[string]interface{}, error) {
	baseKey := rolesKey + ":" + roleID + ":"
	assignations := make(map[string]interface{})

	modules, err := r.c.SMembers(ctx, baseKey+"mo").Result()
	if err != nil {
		return nil, err
	}
	for _, m := range modules {
		module, err := r.moduleStructure(ctx, m)
		if err != nil {
			return nil, err
		}
		module.Access = true

		// get submodules from redis
		key := baseKey + "sm:" + m
		submodules, err := r.c.SMembers(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// get actions from redis
		key = baseKey + "ac:" + m
		for i := range module.SubModules {
			module.SubModules[i].Sections = nil

			actions, err := r.c.SMembers(ctx, key+":"+module.SubModules[i].Name).Result()
			if err != nil {
				return nil, err
			}
			// check against actions from redis
			if r.contains(submodules, module.SubModules[i].Name) {
				module.SubModules[i].Access = true
				for k := range module.SubModules[i].Actions {
					if r.contains(actions, k) {
						action := module.SubModules[i].Actions[k]
						action.Allowed = true
						module.SubModules[i].Actions[k] = action
					}
				}
			} else {
				module.SubModules[i].Actions = nil
			}
		}
		assignations[m] = module
	}
	return assignations, err
}

// RemoveActionsByUser Remove action list for a given user
func (r *actionsRepo) RemoveActionsByUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf(actionsKey, userID)
	_, err := r.c.Del(ctx, key).Result()
	return err
}

// UpdateActionList update the list of actions to quick access while checking permissions
func (r *actionsRepo) UpdateActionList(ctx context.Context, roleID string) error {
	key := fmt.Sprintf(roleUserKey, roleID)
	users, err := r.c.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	actions, err := r.actionsForRole(ctx, roleID)
	for _, userID := range users {
		key = fmt.Sprintf(hasPesmissionKey, userID)
		_, err = r.c.SAdd(ctx, key, actions).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// moduleStructure returns the modules, submodules and sections structure for a given module
func (r *actionsRepo) moduleStructure(ctx context.Context, name string) (*entities.Module, error) {
	key := accessTemplateKey + ":" + name
	var module entities.Module
	j, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(j), &module)
	if err != nil {
		return nil, err
	}
	return &module, err
}

func (r *actionsRepo) actionsForRole(ctx context.Context, roleID string) ([]string, error) {
	key := fmt.Sprintf(roleActionsKeys, roleID)
	actionKeyList, err := r.c.Keys(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var actionList []string
	for _, aKey := range actionKeyList {
		actions, err := r.c.SMembers(ctx, aKey).Result()
		if err != nil {
			return nil, err
		}
		actionList = append(actionList, actions...)
	}
	return actionList, nil

}

func (r *actionsRepo) contains(list []string, el string) bool {
	for _, e := range list {
		if e == el {
			return true
		}
	}
	return false
}
