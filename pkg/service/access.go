package service

import (
	"context"
	"errors"
	"time"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

// AccessService service interface
type AccessService interface {
	// Login log in a user and return access and refresh tokens or an error
	Login(context.Context, string) (*entities.LoggedUser, error)
	// VerifyToken check if a token is valid and the user is logged in
	VerifyToken(context.Context, *entities.Token) (string, error)
	// Refresh refresh a token
	RefreshToken(context.Context, *entities.Token) (*entities.Token, error)
	// Logout log out a user for a given token
	Logout(context.Context, *entities.Token) error
}

type access struct {
	repo   repository.UsersRepository
	config configuration.SecurityConfig
}

// NewAccessService return a new access service instance
func NewAccessService(usersRepo repository.UsersRepository, config configuration.SecurityConfig) AccessService {
	return &access{
		repo:   usersRepo,
		config: config,
	}
}

// Login log in a user by email and return access and refresh tokens or an error
func (ga *access) Login(ctx context.Context, email string) (*entities.LoggedUser, error) {
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
func (ga *access) VerifyToken(ctx context.Context, token *entities.Token) (string, error) {
	claims, err := ga.getTokenClaims(ctx, token.Access)
	if err != nil {
		return "", err
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
func (ga *access) RefreshToken(ctx context.Context, token *entities.Token) (*entities.Token, error) {
	claims, err := ga.getTokenClaims(ctx, token.Access)
	if err == nil {
		accessKey := claims["access_uuid"].(string)
		err = ga.repo.DeleteToken(ctx, accessKey)
		if err != nil {
			return nil, err
		}
	}
	claims, err = ga.getTokenClaims(ctx, token.Refresh)
	if err != nil {
		return nil, err
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
func (ga *access) Logout(ctx context.Context, token *entities.Token) error {
	claims, err := ga.getTokenClaims(ctx, token.Access)
	if err == nil {
		accessKey := claims["access_uuid"].(string)
		err = ga.repo.DeleteToken(ctx, accessKey)
		if err != nil {
			return err
		}
	}
	claims, err = ga.getTokenClaims(ctx, token.Refresh)
	if err == nil {
		refreshKey := claims["refresh_uuid"].(string)
		err = ga.repo.DeleteToken(ctx, refreshKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ga *access) createToken(user *entities.User) (*entities.StoredToken, error) {
	aUUDI := xid.New().String()
	aExp := time.Now().Add(time.Minute * time.Duration(ga.config.JWTTokenExpiration)).Unix()
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["access_uuid"] = aUUDI
	claims["exp"] = aExp
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	atoken, err := t.SignedString([]byte(ga.config.JWTSecret))
	if err != nil {
		return nil, errors.New("Unable to create token")
	}

	rUUDI := xid.New().String()
	rExp := time.Now().Add(time.Minute * time.Duration(ga.config.JWTRefreshExpiration)).Unix()
	claims = jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["refresh_uuid"] = rUUDI
	claims["exp"] = rExp
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rtoken, err := t.SignedString([]byte(ga.config.JWTSecret))
	if err != nil {
		return nil, errors.New("Unable to create token")
	}
	return &entities.StoredToken{
		ID:             user.ID,
		AccessToken:    atoken,
		AccessUUID:     aUUDI,
		AccessExpires:  aExp,
		RefreshToken:   rtoken,
		RefreshUUID:    rUUDI,
		RefreshExpires: rExp,
	}, nil
}

func (ga *access) saveUserToken(ctx context.Context, user *entities.User) (*entities.LoggedUser, error) {
	token, err := ga.createToken(user)
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

func (ga *access) getTokenClaims(ctx context.Context, token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong signed method")
		}
		return []byte(ga.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok && !parsed.Valid {
		return nil, err
	}
	return claims, nil
}
