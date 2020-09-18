package utils

import (
	"errors"
	"time"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
)

// StoredToken stored token struct
type StoredToken struct {
	ID             string
	AccessToken    string
	AccessUUID     string
	AccessExpires  int64
	RefreshToken   string
	RefreshUUID    string
	RefreshExpires int64
}

// JwtHandler interface
type JwtHandler interface {
	CreateToken(ID string) (*StoredToken, error)
	GetTokenClaims(token string) (jwt.MapClaims, error)
}

type jwtHandler struct {
	JWTSecret            string
	JWTTokenExpiration   int
	JWTRefreshExpiration int
}

// NewJwtHandler return a new JWT handler instance
func NewJwtHandler(config configuration.SecurityConfig) JwtHandler {
	return &jwtHandler{
		JWTSecret:            config.JWTSecret,
		JWTTokenExpiration:   config.JWTTokenExpiration,
		JWTRefreshExpiration: config.JWTRefreshExpiration,
	}
}

func (h *jwtHandler) CreateToken(ID string) (*StoredToken, error) {
	aUUDI := xid.New().String()
	aExp := time.Now().Add(time.Minute * time.Duration(h.JWTTokenExpiration)).Unix()
	claims := jwt.MapClaims{}
	claims["user_id"] = ID
	claims["access_uuid"] = aUUDI
	claims["exp"] = aExp
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	atoken, err := t.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return nil, errors.New("Unable to create token")
	}

	rUUDI := xid.New().String()
	rExp := time.Now().Add(time.Minute * time.Duration(h.JWTRefreshExpiration)).Unix()
	claims = jwt.MapClaims{}
	claims["user_id"] = ID
	claims["refresh_uuid"] = rUUDI
	claims["exp"] = rExp
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rtoken, err := t.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return nil, errors.New("Unable to create token")
	}
	return &StoredToken{
		ID:             ID,
		AccessToken:    atoken,
		AccessUUID:     aUUDI,
		AccessExpires:  aExp,
		RefreshToken:   rtoken,
		RefreshUUID:    rUUDI,
		RefreshExpires: rExp,
	}, nil
}

func (h *jwtHandler) GetTokenClaims(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong signed method")
		}
		return []byte(h.JWTSecret), nil
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
