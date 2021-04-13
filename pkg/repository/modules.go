package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"encoding/json"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/go-redis/redis/v8"
)

// ModulesRepository modules repository
type ModulesRepository interface {
	// AssignModule assign module to a role
	AssignModule(ctx context.Context, roleID string, module string) error
	// UnassignModule unassign modules from a role
	UnassignModule(ctx context.Context, roleID string, module string) error
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
	// ModuleStructure returns the modules, submodules and sections structure for a given module
	ModuleStructure(ctx context.Context, name string) (*entities.Module, error)
	// AssignationsByRole get a list of modules, submodules and sections assigned to the role
	AssignationsByRole(ctx context.Context, roleID string) (map[string]interface{}, error)
	// GetAccessList get the modules, submodules and sections assigned to a user
	GetAccessList(ctx context.Context, userID string) (string, error)
	// SetAccessList sets the access list for a given user based on all assigned roles
	SetAccessList(ctx context.Context, userID string) error
	// RemoveAccessByUser Remove access list for a given user
	RemoveAccessByUser(ctx context.Context, userID string) error
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

// AssignModule assign modules to a role
func (r *modulesRepo) AssignModule(ctx context.Context, roleID string, module string) error {
	key := accessTemplateKey + ":" + module
	ok, err := r.c.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if ok == 0 {
		return errors.New("Module not found")
	}
	key = rolesKey + ":" + roleID + ":mo"
	_, err = r.c.SAdd(ctx, key, module).Result()
	return err
}

// UnassignModule unassign modules from a role
func (r *modulesRepo) UnassignModule(ctx context.Context, roleID string, module string) error {
	key := rolesKey + ":" + roleID + ":mo"
	_, err := r.c.SRem(ctx, key, module).Result()
	return err
}

// AssignSubModules assign submodules to a role
func (r *modulesRepo) AssignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	key := fmt.Sprintf(roleSubModulesKey, rolesKey, roleID, module)
	_, err := r.c.SAdd(ctx, key, submodules).Result()
	return err
}

// UnassignSubModules unassign submodules from a role
func (r *modulesRepo) UnassignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	key := fmt.Sprintf(roleSubModulesKey, rolesKey, roleID, module)
	_, err := r.c.SRem(ctx, key, submodules).Result()
	return err
}

// AssignSections assign sections to a role
func (r *modulesRepo) AssignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	key := fmt.Sprintf(roleSectionsKey, rolesKey, roleID, module, submodule)
	_, err := r.c.SAdd(ctx, key, sections).Result()
	return err
}

// UnassignSections unassign sections from a role
func (r *modulesRepo) UnassignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	key := fmt.Sprintf(roleSectionsKey, rolesKey, roleID, module, submodule)
	_, err := r.c.SRem(ctx, key, sections).Result()
	return err
}

// ModulesList returns a list of available modules
func (r *modulesRepo) ModulesList(ctx context.Context) ([]string, error) {
	keys, err := r.c.Keys(ctx, accessTemplateKey+"*").Result()
	if err != nil {
		return nil, err
	}
	var modules []string
	for _, k := range keys {
		modules = append(modules, strings.Replace(k, accessTemplateKey+":", "", 1))
	}
	return modules, nil
}

// ModuleStructure returns the modules, submodules and sections structure for a given module
func (r *modulesRepo) ModuleStructure(ctx context.Context, name string) (*entities.Module, error) {
	key := accessTemplateKey + ":" + name
	var module entities.Module
	j, err := r.c.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}
	err = json.Unmarshal([]byte(j), &module)
	if err != nil {
		return nil, err
	}
	return &module, err
}

// AssignationsByRole get a list of modules, submodules and sections assigned to the role
func (r *modulesRepo) AssignationsByRole(ctx context.Context, roleID string) (map[string]interface{}, error) {
	baseKey := rolesKey + ":" + roleID + ":"
	assignations := make(map[string]interface{})

	modules, err := r.c.SMembers(ctx, baseKey+"mo").Result()
	if err != nil {
		return nil, err
	}
	for _, m := range modules {
		module, err := r.ModuleStructure(ctx, m)
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
		// get sections from redis
		key = baseKey + "se:" + m
		for i := range module.SubModules {
			module.SubModules[i].Actions = nil
			sections, err := r.c.SMembers(ctx, key+":"+module.SubModules[i].Name).Result()
			if err != nil {
				return nil, err
			}
			// check against submodules from redis
			if r.contains(submodules, module.SubModules[i].Name) {
				module.SubModules[i].Access = true
				for k := range module.SubModules[i].Sections {
					if r.contains(sections, k) {
						module.SubModules[i].Sections[k] = true
					}
				}
			}
		}
		assignations[m] = module
	}
	return assignations, err
}

// GetAccessList get the modules, submodules and sections assigned to a user
// TODO: Move this logic to subscriber in order to have the final json stored and updated by userID
func (r *modulesRepo) GetAccessList(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf(accessKey, userID)
	j, err := r.c.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if j == "" {
		return "", errors.New("The user has no access list defined")
	}
	return j, nil
}

// SetAccessList sets the access list for a given user based on all assigned roles
func (r *modulesRepo) SetAccessList(ctx context.Context, userID string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	roles, err := r.c.SMembers(ctx, key).Result()
	if err != nil {
		return err
	}
	sort.Strings(roles)
	assignations := make(map[string]interface{})
	for _, role := range roles {
		assignedModules, err := r.AssignationsByRole(ctx, role)
		if err != nil {
			return err
		}
		// Modules are override with the latest role
		for module, assignation := range assignedModules {
			assignations[module] = assignation
		}
	}
	if err != nil && err != redis.Nil {
		return err
	}
	j, err := json.Marshal(assignations)
	key = fmt.Sprintf(accessKey, userID)
	_, err = r.c.Set(ctx, key, j, 0).Result()
	return err
}

// RemoveAccessByUser Remove access list for a given user
func (r *modulesRepo) RemoveAccessByUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf(accessKey, userID)
	_, err := r.c.Del(ctx, key).Result()
	return err
}

func (r *modulesRepo) contains(list []string, el string) bool {
	for _, e := range list {
		if e == el {
			return true
		}
	}
	return false
}
