package repository

import (
	"context"
	"encoding/json"

	"github.com/StevenRojas/goaccess/pkg/entities"

	"github.com/go-redis/redis/v8"
)

// InitRepository interface
type InitRepository interface {
	// IsSetConfig checks if the DB is initialized
	IsSetConfig(ctx context.Context) (bool, error)
	// SetAsConfigured sets the DB as initialized
	SetAsConfigured(ctx context.Context) error
	// UnsetConfig sets the DB as not initialized
	UnsetConfig(ctx context.Context) error
	// AddModule add a module in the DB
	AddModule(ctx context.Context, module entities.ModuleInit) error
}

type initRepo struct {
	c *redis.Client
}

// NewInitRepository creates a new repository instance
func NewInitRepository(ctx context.Context, client *redis.Client) (InitRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &initRepo{
		c: client,
	}, nil
}

// IsSetConfig checks if the DB is initialized
func (r *initRepo) IsSetConfig(ctx context.Context) (bool, error) {
	isset, err := r.c.Get(ctx, isSetKey).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	return isset == "true", nil
}

// SetAsConfigured sets the DB as initialized
func (r *initRepo) SetAsConfigured(ctx context.Context) error {
	_, err := r.c.Set(ctx, isSetKey, "true", 0).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// UnsetConfig sets the DB as not initialized
func (r *initRepo) UnsetConfig(ctx context.Context) error {
	pipe := r.c.Pipeline()
	iter := r.c.Scan(ctx, 0, configKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AddModule add a module in the DB
func (r *initRepo) AddModule(ctx context.Context, module entities.ModuleInit) error {
	key := accessTemplateKey + ":" + module.Name
	j, _ := json.Marshal(module)
	_, err := r.c.Set(ctx, key, j, -1).Result()
	return err
}
