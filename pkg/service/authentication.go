package service

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess/pkg/utils"

	"github.com/StevenRojas/goaccess/pkg/repository"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

// AuthenticationService service interface
type AuthenticationService interface {
	// Register a user
	Register(context.Context, *entities.User) error
	// Unregister a user
	Unregister(context.Context, *entities.User) error
	// Login log in a user and return access and refresh tokens or an error
	Login(context.Context, string) (*entities.LoggedUser, error)
	// VerifyToken check if a token is valid and the user is logged in
	VerifyToken(context.Context, string) (string, error)
	// Refresh refresh a token
	RefreshToken(context.Context, string) (*entities.Token, error)
	// Logout log out a user for a given token
	Logout(context.Context, *entities.Token) error
}

type authentication struct {
	repo       repository.UsersRepository
	jwtHandler utils.JwtHandler
}

// NewAuthenticationService return a new authentication service instance
func NewAuthenticationService(usersRepo repository.UsersRepository, jwtHandler utils.JwtHandler) AuthenticationService {
	return &authentication{
		repo:       usersRepo,
		jwtHandler: jwtHandler,
	}
}

// Register a user
func (ga *authentication) Register(ctx context.Context, user *entities.User) error {
	return ga.repo.Register(ctx, user)
}

// Unregister a user
func (ga *authentication) Unregister(ctx context.Context, user *entities.User) error {
	return ga.repo.Unregister(ctx, user)
}

// Login log in a user by email and return access and refresh tokens or an error
func (ga *authentication) Login(ctx context.Context, email string) (*entities.LoggedUser, error) {
	user, err := ga.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found")
	}
	return ga.saveUserToken(ctx, user)
}

// VerifyToken check if a token is valid and the user is logged in
func (ga *authentication) VerifyToken(ctx context.Context, token string) (string, error) {
	claims, err := ga.jwtHandler.GetTokenClaims(token)
	if err != nil {
		return "", err
	}
	if _, ok := claims["access_uuid"]; !ok {
		return "", errors.New("Invalid token")
	}
	if _, ok := claims["user_id"]; !ok {
		return "", errors.New("Invalid token")
	}
	accessKey := claims["access_uuid"].(string)
	user, err := ga.repo.GetUserByToken(ctx, accessKey)
	if err != nil {
		return "", err
	}
	if user == nil || user.ID != claims["user_id"].(string) {
		return "", errors.New("Invalid or expired token")
	}
	return user.ID, nil
}

// Refresh refresh a token
func (ga *authentication) RefreshToken(ctx context.Context, token string) (*entities.Token, error) {
	claims, err := ga.jwtHandler.GetTokenClaims(token)
	if err != nil {
		return nil, err
	}
	if _, ok := claims["refresh_uuid"]; !ok {
		return nil, errors.New("Invalid token")
	}
	if _, ok := claims["user_id"]; !ok {
		return nil, errors.New("Invalid token")
	}
	refreshKey := claims["refresh_uuid"].(string)
	user, err := ga.repo.GetUserByToken(ctx, refreshKey)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID != claims["user_id"].(string) {
		return nil, errors.New("Invalid or expired token")
	}
	err = ga.repo.DeleteToken(ctx, refreshKey)
	if err != nil {
		return nil, err
	}
	loggedUser, err := ga.saveUserToken(ctx, user)
	return loggedUser.Token, err
}

// Logout log out a user for a given token
func (ga *authentication) Logout(ctx context.Context, token *entities.Token) error {
	claims, err := ga.jwtHandler.GetTokenClaims(token.Access)
	if err == nil {
		accessKey := claims["access_uuid"].(string)
		err = ga.repo.DeleteToken(ctx, accessKey)
		if err != nil {
			return err
		}
	}
	claims, err = ga.jwtHandler.GetTokenClaims(token.Refresh)
	if err == nil {
		refreshKey := claims["refresh_uuid"].(string)
		err = ga.repo.DeleteToken(ctx, refreshKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ga *authentication) saveUserToken(ctx context.Context, user *entities.User) (*entities.LoggedUser, error) {
	token, err := ga.jwtHandler.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}
	err = ga.repo.StoreTokens(ctx, token)
	if err != nil {
		return nil, err
	}
	return &entities.LoggedUser{
		User: user,
		Token: &entities.Token{
			Access:  token.AccessToken,
			Refresh: token.RefreshToken,
		},
	}, err
}
