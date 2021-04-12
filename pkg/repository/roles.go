package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// RolesRepository roles repository
type RolesRepository interface {
	// AddRole add a role and return its ID
	AddRole(ctx context.Context, name string) (string, error)
	//EditRole edit the role name
	EditRole(ctx context.Context, ID string, name string) error
	// DeleteRole removes a role and its relation with users
	DeleteRole(ctx context.Context, ID string) error
	// IsValidRole check if a role exist
	IsValidRole(ctx context.Context, ID string) (bool, error)
	// AssingRole assign role to a user
	AssignRole(ctx context.Context, userID string, roleID string) error
	// UnassignRole unassign role from a user
	UnassignRole(ctx context.Context, userID string, roleID string) error
	// UsersByRole get a list of users assigned to a given role
	UsersByRole(ctx context.Context, roleID string) ([]string, error)
	// GetRoles get a list of all roles
	GetRoles(ctx context.Context) (map[string]string, error)
	// RolesByUser get a list of roles assigned to a user
	RolesByUser(ctx context.Context, userID string) (map[string]string, error)
}

type roleRepo struct {
	c *redis.Client
}

// NewRolesRepository creates a new repository instance
func NewRolesRepository(ctx context.Context, client *redis.Client) (RolesRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &roleRepo{
		c: client,
	}, nil
}

// AddRole add a role and return its ID
func (r *roleRepo) AddRole(ctx context.Context, name string) (string, error) {
	id, err := r.c.Incr(ctx, roleIDKey).Result()
	if err != nil {
		return "", err
	}
	rid := "r" + strconv.FormatInt(id, 10)
	_, err = r.c.HMSet(ctx, rolesKey, rid, name).Result()
	if err != nil {
		return "", err
	}
	return rid, nil
}

//EditRole edit the role name
func (r *roleRepo) EditRole(ctx context.Context, ID string, name string) error {
	_, err := r.c.HMSet(ctx, rolesKey, ID, name).Result()
	return err
}

// DeleteRole removes a role and its relation with users
func (r *roleRepo) DeleteRole(ctx context.Context, ID string) error {
	branchKeys, _ := r.c.Keys(ctx, rolesKey+":"+ID+":*").Result()
	users, err := r.UsersByRole(ctx, ID)
	if err != nil {
		return err
	}
	pipe := r.c.Pipeline()
	for _, userID := range users {
		key := fmt.Sprintf(userRoleKey, userID)
		pipe.SRem(ctx, key, ID) // remove role member from userrole:1
	}
	key := fmt.Sprintf(roleUserKey, ID)
	pipe.Del(ctx, key)
	pipe.HDel(ctx, rolesKey, ID).Result()
	for _, k := range branchKeys {
		pipe.Del(ctx, k)
	}
	// FIXME: Delete access:user and actions:user if the deleted role is the only one assigned to the user
	pipe.Del(ctx, rolesKey+":"+ID)
	_, err = pipe.Exec(ctx)
	return err
}

// IsValidRole check if a role exist
func (r *roleRepo) IsValidRole(ctx context.Context, ID string) (bool, error) {
	return r.c.HExists(ctx, rolesKey, ID).Result()
}

// AssingRole assign role to a user
func (r *roleRepo) AssignRole(ctx context.Context, userID string, roleID string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	pipe := r.c.Pipeline()
	pipe.SAdd(ctx, key, roleID) //userrole
	key = fmt.Sprintf(roleUserKey, roleID)
	pipe.SAdd(ctx, key, userID) // roleuser
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// UnassignRole unassign role from a user
func (r *roleRepo) UnassignRole(ctx context.Context, userID string, roleID string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	pipe := r.c.Pipeline()
	pipe.SRem(ctx, key, roleID)
	key = fmt.Sprintf(roleUserKey, roleID)
	pipe.SRem(ctx, key, userID)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// UsersByRole get a list of users assigned to a given role
func (r *roleRepo) UsersByRole(ctx context.Context, roleID string) ([]string, error) {
	key := fmt.Sprintf(roleUserKey, roleID)
	users, err := r.c.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// RolesByUser get a list of roles assigned to a user
func (r *roleRepo) RolesByUser(ctx context.Context, userID string) (map[string]string, error) {
	roleList, err := r.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf(userRoleKey, userID)
	roleKeys, err := r.c.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	roles := map[string]string{}
	for _, key := range roleKeys {
		if role, ok := roleList[key]; ok {
			roles[key] = role
		}
	}
	return roles, nil
}

// GetRoles get a list of all roles
func (r *roleRepo) GetRoles(ctx context.Context) (map[string]string, error) {
	roles, err := r.c.HGetAll(ctx, "roles").Result()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
