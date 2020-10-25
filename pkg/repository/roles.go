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
	// AssingRoles assign roles to a user
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	// UnassignRoles unassign roles from a user
	UnassignRoles(ctx context.Context, userID string, roleIDs []string) error
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
	_, err := r.c.HDel(ctx, rolesKey, ID).Result()
	return err
}

// AssingRoles assign roles to a user
func (r *roleRepo) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	pipe := r.c.Pipeline()
	pipe.SAdd(ctx, key, roleIDs) //userrole
	for _, roleID := range roleIDs {
		key = fmt.Sprintf(roleUserKey, roleID)
		pipe.SAdd(ctx, key, userID) // roleuser
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// UnassignRoles unassign roles from a user
func (r *roleRepo) UnassignRoles(ctx context.Context, userID string, roleIDs []string) error {
	key := fmt.Sprintf(userRoleKey, userID)
	pipe := r.c.Pipeline()
	pipe.SRem(ctx, key, roleIDs)
	for _, roleID := range roleIDs {
		key = fmt.Sprintf(roleUserKey, roleID)
		pipe.SRem(ctx, key, userID)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
