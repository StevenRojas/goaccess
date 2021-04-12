package repository

import (
	"context"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/utils"
	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/mock"
)

// UsersRepositoryMock interface
type UsersRepositoryMock interface {
	// Register a user
	Register(context.Context, *entities.User) error
	// Unregister a user
	Unregister(context.Context, *entities.User) error
	// GetUsers get a user by ID
	GetUsers(context.Context) ([]entities.User, error)
	// GetUsersByRole get a user by ID
	GetUsersByRole(context.Context, string) ([]entities.User, error)
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

// UsersRepoMock users repo mock
type UsersRepoMock struct {
	M mock.Mock
}

// NewUsersRepositoryMock creates a new repository instance
func NewUsersRepositoryMock(ctx context.Context, client *redis.Client) (UsersRepository, error) {
	return new(UsersRepoMock), nil
}

// GetUsers get user list
func (r *UsersRepoMock) GetUsers(ctx context.Context) ([]entities.User, error) {
	args := r.M.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

// GetUsersByRole get user list by role
func (r *UsersRepoMock) GetUsersByRole(ctx context.Context, roleID string) ([]entities.User, error) {
	args := r.M.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

// Register a user
func (r *UsersRepoMock) Register(ctx context.Context, user *entities.User) error {
	args := r.M.Called(user.ID)
	return args.Error(1)
}

// Unregister a user
func (r *UsersRepoMock) Unregister(ctx context.Context, user *entities.User) error {
	args := r.M.Called(user.ID)
	return args.Error(1)
}

// GetUserByID get a user by ID
func (r *UsersRepoMock) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := r.M.Called(id)
	return args.Get(0).(*entities.User), args.Error(1)
}

// GetUserByEmail get a user by email
func (r *UsersRepoMock) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := r.M.Called(email)
	user := args.Get(0)
	if user != nil {
		return args.Get(0).(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetUserByToken get a user by email
func (r *UsersRepoMock) GetUserByToken(ctx context.Context, token string) (*entities.User, error) {
	args := r.M.Called(token)
	return args.Get(0).(*entities.User), args.Error(1)
}

// StoreTokens get a user by email
func (r *UsersRepoMock) StoreTokens(ctx context.Context, token *utils.StoredToken) error {
	args := r.M.Called(token)
	return args.Error(0)
}

// DeleteToken get a user by email
func (r *UsersRepoMock) DeleteToken(ctx context.Context, key string) error {
	args := r.M.Called(key)
	return args.Error(0)
}

// IsValidUser check if a user exist
func (r *UsersRepoMock) IsValidUser(ctx context.Context, id string) (bool, error) {
	args := r.M.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
