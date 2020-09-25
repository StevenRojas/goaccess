package repository

import (
	"context"

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
	return "", nil
}

//EditRole edit the role name
func (r *roleRepo) EditRole(ctx context.Context, ID string, name string) error {
	return nil
}

// DeleteRole removes a role and its relation with users
func (r *roleRepo) DeleteRole(ctx context.Context, ID string) error {
	return nil
}
