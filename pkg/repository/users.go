package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/redis"
)

// UsersRepository interface
type UsersRepository interface {
	GetUserByID(context.Context, string) (*entities.User, error)
	GetUserByEmail(context.Context, string) (*entities.User, error)
	GetUserByToken(context.Context, string) (*entities.User, error)
	StoreTokens(context.Context, *entities.StoredToken) error
	DeleteToken(context.Context, string) error
}

type repo struct {
	c redis.RedisClient
}

// NewUsersRepository creates a new repository instance
func NewUsersRepository(ctx context.Context, client redis.RedisClient) (UsersRepository, error) {
	return &repo{
		c: client,
	}, nil
}

// GetUserByID get a user by ID
func (r *repo) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	key := fmt.Sprintf("user:%s", id)
	result, err := r.c.GetAttributes(ctx, key)
	if err != nil {
		return nil, err
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
	id, err := r.c.GetHash(ctx, "users", email)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, id)
}

// getUserByToken get a user by token
func (r *repo) GetUserByToken(ctx context.Context, token string) (*entities.User, error) {
	id, err := r.c.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, id)
}

// StoreTokens store access and refresh token hashes with an expiration period
func (r *repo) StoreTokens(ctx context.Context, token *entities.StoredToken) error {
	at := time.Unix(token.AccessExpires, 0)
	rt := time.Unix(token.RefreshExpires, 0)
	now := time.Now()
	err := r.c.Set(ctx, token.AccessUUID, token.ID, at.Sub(now))
	if err != nil {
		return err
	}
	err = r.c.Set(ctx, token.RefreshUUID, token.ID, rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken delete token key
func (r *repo) DeleteToken(ctx context.Context, key string) error {
	err := r.c.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
