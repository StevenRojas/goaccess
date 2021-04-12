package service

import (
	"context"
	"testing"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"
	"github.com/StevenRojas/goaccess/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceSuite struct {
	svc  AuthenticationService
	repo *repository.UsersRepoMock
	suite.Suite
}

func (s *serviceSuite) SetupTest() {
	//ctx := context.TODO()
	s.repo = new(repository.UsersRepoMock)
	jwtHander := utils.NewJwtHandlerMock(configuration.SecurityConfig{
		JWTSecret:            "secret!",
		JWTTokenExpiration:   10,
		JWTRefreshExpiration: 20,
	})
	s.svc = NewAuthenticationService(s.repo, jwtHander)
}

func TestAccessService(t *testing.T) {
	ss := new(serviceSuite)
	suite.Run(t, ss)
}

func (s *serviceSuite) TestLoginNotRegistered() {
	t := s.T()
	email := "notRegistered@gmail.com"
	s.repo.M.On("GetUserByEmail", email).Return(nil, nil)
	_, err := s.svc.Login(context.TODO(), email)
	assert.NotNil(t, err)
	assert.Equal(t, "User not found", err.Error())
}

func (s *serviceSuite) TestLoginSuccess() {
	t := s.T()
	email := "srojas@gmail.com"
	expected := &entities.User{
		ID:    "1",
		Email: email,
		Name:  "steven rojas",
	}
	s.repo.M.On("GetUserByEmail", email).Return(expected, nil)
	s.repo.M.On("StoreTokens", &utils.StoredToken{
		ID:             "1",
		AccessToken:    "a_jwt",
		AccessUUID:     "a_uuid",
		AccessExpires:  10,
		RefreshToken:   "r_jwt",
		RefreshUUID:    "r_uuid",
		RefreshExpires: 20,
	}).Return(nil)
	token, err := s.svc.Login(context.TODO(), email)
	assert.Nil(t, err)
	assert.Equal(t, expected, token.User)
	assert.NotNil(t, token.Token)
}

// func (s *serviceSuite) TestValidToken() {
// 	t := s.T()
// 	token := "a_jwt"
// 	eUser := &entities.User{
// 		ID:    "1",
// 		Email: "srojas@gamil.com",
// 		Name:  "steven rojas",
// 	}
// 	s.repo.M.On("GetTokenClaims", token).Return(nil, nil)
// 	key := "a_uuid"
// 	s.repo.M.On("GetUserByToken", key).Return(eUser, nil)
// 	ID, err := s.svc.VerifyToken(context.TODO(), &entities.Token{
// 		Access:  "a_jwt",
// 		Refresh: "r_jwt",
// 	})
// 	assert.Nil(t, err)
// 	assert.Equal(t, eUser.ID, ID)
// }
