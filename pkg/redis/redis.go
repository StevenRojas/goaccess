package redis

import (
	"context"
	"errors"
	"time"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/go-redis/redis/v8"
)

// RedisClient interface
type RedisClient interface {
	Ping() error
	GetHash(context.Context, string, string) (string, error)
	GetAttributes(context.Context, string) (map[string]string, error)
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error
	Del(context.Context, string) error
}

type redisClient struct {
	c *redis.Client
}

// NewRedisClient creates a new redis client instance
func NewRedisClient(ctx context.Context, config configuration.RedisConfig) (RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pass,
		DB:       config.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &redisClient{
		c: rdb,
	}, nil
}

func (rc *redisClient) GetHash(ctx context.Context, key string, value string) (string, error) {
	id, err := rc.c.HGet(ctx, key, value).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if id == "" {
		return "", errors.New("Not found")
	}
	return id, nil
}

func (rc *redisClient) GetAttributes(ctx context.Context, key string) (map[string]string, error) {
	att, err := rc.c.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if att == nil {
		return nil, errors.New("Not found")
	}
	return att, nil
}

func (rc *redisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.c.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if value == "" {
		return "", errors.New("Not found")
	}
	return value, nil
}

func (rc *redisClient) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	_, err := rc.c.Set(ctx, key, value, exp).Result()
	return err
}

func (rc *redisClient) Del(ctx context.Context, key string) error {
	_, err := rc.c.Del(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// Ping ping to server
func (rc *redisClient) Ping() error {
	_, err := rc.c.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	return nil
}
