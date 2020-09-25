package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const isSetKey = "isset"
const configKey = "c:"
const actionsKey = "ac:"
const modulesKey = "c:mo"
const subModulesPrefix = "c:sm:"
const sectionsPrefix = "c:se:"
const actionsPrefix = "c:ac:"

// InitRepository interface
type InitRepository interface {
	IsSetConfig(ctx context.Context) (bool, error)
	SetAsConfigured(ctx context.Context) error
	UnsetConfig(ctx context.Context) error
	SetModules(ctx context.Context, modules []string) error
	SetSubModules(ctx context.Context, module string, submodules []string) error
	SetSection(ctx context.Context, module string, submodule string, sections []string) error
	SetActions(ctx context.Context, module string, submodule string, actions []string) error
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

func (r *initRepo) IsSetConfig(ctx context.Context) (bool, error) {
	isset, err := r.c.Get(ctx, isSetKey).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	return isset == "true", nil
}

func (r *initRepo) SetAsConfigured(ctx context.Context) error {
	_, err := r.c.Set(ctx, isSetKey, "true", 0).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

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

func (r *initRepo) SetModules(ctx context.Context, modules []string) error {
	_, err := r.c.LPush(ctx, modulesKey, modules).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *initRepo) SetSubModules(ctx context.Context, module string, submodules []string) error {
	key := subModulesPrefix + module
	_, err := r.c.LPush(ctx, key, submodules).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *initRepo) SetSection(ctx context.Context, module string, submodule string, sections []string) error {
	key := sectionsPrefix + module[:2] + ":" + submodule
	_, err := r.c.LPush(ctx, key, sections).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *initRepo) SetActions(ctx context.Context, module string, submodule string, actions []string) error {
	key := actionsPrefix + module[:2] + ":" + submodule
	_, err := r.c.LPush(ctx, key, actions).Result()
	if err != nil {
		return err
	}
	return nil
}
