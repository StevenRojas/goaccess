package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/utils"
	"github.com/go-redis/redis/v8"
)

// UsersRepository interface
type UsersRepository interface {
	// GetUserByID get a user by ID
	GetUserByID(context.Context, string) (*entities.User, error)
	// GetUserByEmail get a user by email
	GetUserByEmail(context.Context, string) (*entities.User, error)
	// GetUserByToken get a user by token
	GetUserByToken(context.Context, string) (*entities.User, error)
	// StoreTokens store access and refresh token hashes with an expiration period
	StoreTokens(context.Context, *utils.StoredToken) error
	// DeleteToken delete token key
	DeleteToken(context.Context, string) error
}

type repo struct {
	c *redis.Client
}

// NewUsersRepository creates a new repository instance
func NewUsersRepository(ctx context.Context, client *redis.Client) (UsersRepository, error) {
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &repo{
		c: client,
	}, nil
}

// GetUserByID get a user by ID
func (r *repo) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	key := fmt.Sprintf("user:%s", id)
	result, err := r.c.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("Not found")
	}
	user := &entities.User{
		ID:    result["id"],
		Email: result["email"],
		Name:  result["name"],
	}
	return user, nil
}

// GetUserByEmail get a user by email
func (r *repo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	id, err := r.c.HGet(ctx, "users", email).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if id == "" {
		return nil, errors.New("Not found")
	}
	return r.GetUserByID(ctx, id)
}

// GetUserByToken get a user by token
func (r *repo) GetUserByToken(ctx context.Context, token string) (*entities.User, error) {
	id, err := r.c.Get(ctx, token).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if id == "" {
		return nil, errors.New("Not found")
	}
	return r.GetUserByID(ctx, id)
}

// StoreTokens store access and refresh token hashes with an expiration period
func (r *repo) StoreTokens(ctx context.Context, token *utils.StoredToken) error {
	at := time.Unix(token.AccessExpires, 0)
	rt := time.Unix(token.RefreshExpires, 0)
	now := time.Now()
	_, err := r.c.Set(ctx, token.AccessUUID, token.ID, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	_, err = r.c.Set(ctx, token.RefreshUUID, token.ID, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken delete token key
func (r *repo) DeleteToken(ctx context.Context, key string) error {
	_, err := r.c.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
