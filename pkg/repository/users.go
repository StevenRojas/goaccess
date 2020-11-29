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
	// Register a user
	Register(context.Context, *entities.User) error
	// Unregister a user
	Unregister(context.Context, *entities.User) error
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
	// IsValidUser check if a user exist
	IsValidUser(context.Context, string) (bool, error)
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

// Register a user
func (r *repo) Register(ctx context.Context, user *entities.User) error {
	key := fmt.Sprintf(userKey, user.ID)
	_, err := r.c.HMSet(ctx, key,
		"id", user.ID,
		"email", user.Email,
		"name", user.Name,
		"admin", user.IsAdmin,
	).Result()
	if err != nil {
		return err
	}
	_, err = r.c.HSet(ctx, usersKey, user.Email, user.ID).Result()
	if err != nil {
		return err
	}
	return nil
}

// Unregister a user
func (r *repo) Unregister(ctx context.Context, user *entities.User) error {
	pipe := r.c.Pipeline()
	key := fmt.Sprintf(userKey, user.ID)
	pipe.Del(ctx, key)
	pipe.HDel(ctx, usersKey, user.Email)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID get a user by ID
func (r *repo) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	key := fmt.Sprintf(userKey, id)
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
	id, err := r.c.HGet(ctx, usersKey, email).Result()
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
	key := "tokens:" + token
	id, err := r.c.Get(ctx, key).Result()
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
	key := "tokens:" + token.AccessUUID
	_, err := r.c.Set(ctx, key, token.ID, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	key = "tokens:" + token.RefreshUUID
	_, err = r.c.Set(ctx, key, token.ID, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken delete token key
func (r *repo) DeleteToken(ctx context.Context, key string) error {
	key = "tokens:" + key
	_, err := r.c.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// IsValidUser check if a user exist
func (r *repo) IsValidUser(ctx context.Context, ID string) (bool, error) {
	key := fmt.Sprintf(userKey, ID)
	res, err := r.c.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
