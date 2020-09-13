package service

import (
	"context"
	"testing"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceSuite struct {
	svc AccessService
	suite.Suite
}

func (s *serviceSuite) SetupTest() {
	ctx := context.TODO()
	client, err := redis.NewRedisClient(ctx, configuration.RedisConfig{
		Addr: "localhost:6379",
		Pass: "",
		DB:   10,
	})
	if err != nil {
		panic(err)
	}
	usersRepo, err := repository.NewUsersRepository(ctx, client)
	s.svc = NewAccessService(usersRepo, configuration.SecurityConfig{
		JWTSecret:            "secret!",
		JWTTokenExpiration:   10,
		JWTRefreshExpiration: 20,
	})
	if err != nil {
		panic(err)
	}
}

func TestAccessService(t *testing.T) {
	ss := new(serviceSuite)
	suite.Run(t, ss)
}

func (s *serviceSuite) TestLoginNotRegistered() {
	t := s.T()
	email := "notRegistered@gmail.com"
	_, err := s.svc.Login(context.TODO(), email)
	assert.NotNil(t, err)
	assert.Equal(t, "Not found", err.Error())
}

func (s *serviceSuite) TestLoginSuccess() {
	t := s.T()
	email := "srojas@gmail.com"
	expected := &entities.User{
		ID:    "1",
		Email: email,
		Name:  "steven rojas",
	}
	token, err := s.svc.Login(context.TODO(), email)
	assert.Nil(t, err)
	assert.Equal(t, expected, token.User)
	assert.NotNil(t, token.Token)
}
