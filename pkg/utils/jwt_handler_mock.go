package utils

import (
	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/dgrijalva/jwt-go"
)

// JwtHandlerMock interface
type JwtHandlerMock interface {
	CreateToken(ID string) (*StoredToken, error)
	GetTokenClaims(token string) (jwt.MapClaims, error)
}

type jwtHandlerMock struct {
}

// NewJwtHandlerMock return a new JWT handler instance
func NewJwtHandlerMock(config configuration.SecurityConfig) JwtHandlerMock {
	return &jwtHandlerMock{}
}

func (h *jwtHandlerMock) CreateToken(ID string) (*StoredToken, error) {
	return &StoredToken{
		ID:             "1",
		AccessToken:    "a_jwt",
		AccessUUID:     "a_uuid",
		AccessExpires:  10,
		RefreshToken:   "r_jwt",
		RefreshUUID:    "r_uuid",
		RefreshExpires: 20,
	}, nil
}

func (h *jwtHandlerMock) GetTokenClaims(token string) (jwt.MapClaims, error) {
	claims := make(map[string]interface{})
	claims["access_uuid"] = "a_uuid"
	claims["user_id"] = "1"
	claims["exp"] = "10"
	return claims, nil
}
