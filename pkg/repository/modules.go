package repository

import (
	"context"
	"errors"
	"fmt"

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
	// ModulesStructure returns a list of available modules for create a new role
	ModulesStructure(ctx context.Context) ([]entities.Module, error)
	// ModuleStructure returns the modules, submodules and sections structure for a given module
	ModuleStructure(ctx context.Context, name string) (*entities.Module, error)
	// AssignationsByRole get a list of modules, submodules and sections assigned to the role
	AssignationsByRole(ctx context.Context, roleID string) (map[string]interface{}, error)
	// GetAccessList get the modules, submodules and sections assigned to a user
	GetAccessList(ctx context.Context, userID string) (string, error)
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
	key := rolesKey + ":" + roleID + ":access:" + module[:2]
	ok, err := r.c.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if ok == 1 {
		return nil
	}
	tKey := accessTemplateKey + ":" + module
	template, err := r.c.Get(ctx, tKey).Result()
	if err == redis.Nil {
		return errors.New("Invalid module name")
	}
	if err != nil {
		return err
	}
	var m entities.Module
	err = json.Unmarshal([]byte(template), &m)
	if err != nil {
		return err
	}
	m.Access = true
	j, _ := json.Marshal(m)
	_, err = r.c.Set(ctx, key, string(j), 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// UnassignModule unassign modules from a role
func (r *modulesRepo) UnassignModule(ctx context.Context, roleID string, module string) error {
	key := rolesKey + ":" + roleID + ":access:" + module[:2]
	_, err := r.c.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// AssignSubModules assign submodules to a role
func (r *modulesRepo) AssignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	key := fmt.Sprintf(roleSubModulesKey, rolesKey, roleID, module[:2])
	_, err := r.c.SAdd(ctx, key, submodules).Result()
	return err
}

// UnassignSubModules unassign submodules from a role
func (r *modulesRepo) UnassignSubModules(ctx context.Context, roleID string, module string, submodules []string) error {
	key := fmt.Sprintf(roleSubModulesKey, rolesKey, roleID, module[:2])
	_, err := r.c.SRem(ctx, key, submodules).Result()
	return err
}

// AssignSections assign sections to a role
func (r *modulesRepo) AssignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	key := fmt.Sprintf(roleSectionsKey, rolesKey, roleID, module[:2], submodule)
	_, err := r.c.SAdd(ctx, key, sections).Result()
	return err
}

// UnassignSections unassign sections from a role
func (r *modulesRepo) UnassignSections(ctx context.Context, roleID string, module string, submodule string, sections []string) error {
	key := fmt.Sprintf(roleSectionsKey, rolesKey, roleID, module[:2], submodule)
	_, err := r.c.SRem(ctx, key, sections).Result()
	return err
}

// ModulesStructure returns a list of available modules, submodules and sections
func (r *modulesRepo) ModulesStructure(ctx context.Context) ([]entities.Module, error) {
	keys, err := r.c.Keys(ctx, accessTemplateKey+"*").Result()
	if err != nil {
		return nil, err
	}
	var modules []entities.Module
	var module entities.Module
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

// ModuleStructure returns the modules, submodules and sections structure for a given module
func (r *modulesRepo) ModuleStructure(ctx context.Context, name string) (*entities.Module, error) {
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
		key := baseKey + "sm:" + m[:2]
		submodules, err := r.c.SMembers(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// get sections from redis
		key = baseKey + "se:" + m[:2]
		sections, err := r.c.SMembers(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		for i := range module.SubModules {
			module.SubModules[i].Actions = nil
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
	j, _ := json.Marshal(assignations)
	fmt.Printf("%v\n", string(j))
	return assignations, err
}

// GetAccessList get the modules, submodules and sections assigned to a user
func (r *modulesRepo) GetAccessList(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf(userRoleKey, userID)
	roles, err := r.c.SMembers(ctx, key).Result()
	fmt.Printf("%v\n", roles)
	// for r := range roles {
	// r.AssignationsByRole(ctx, r)
	// }
	if err != nil && err != redis.Nil {
		return "", err
	}
	return "", nil
}

func (r *modulesRepo) contains(list []string, el string) bool {
	for _, e := range list {
		if e == el {
			return true
		}
	}
	return false
}
