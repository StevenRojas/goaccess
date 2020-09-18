package main

import (
	"context"
	"log"

	"github.com/StevenRojas/goaccess/pkg/utils"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"
	"github.com/go-redis/redis/v8"

	"github.com/StevenRojas/goaccess/pkg/service"

	"github.com/StevenRojas/goaccess/pkg/configuration"
)

func main() {
	ctx := context.Background()
	serviceConfig, err := configuration.Read()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     serviceConfig.Redis.Addr,
		Password: serviceConfig.Redis.Pass,
		DB:       serviceConfig.Redis.DB,
	})

	usersRepo, err := repository.NewUsersRepository(ctx, redisClient)
	if err != nil {
		panic(err)
	}
	jwtHander := utils.NewJwtHandler(serviceConfig.Security)
	s := service.NewAuthenticationService(usersRepo, jwtHander)

	// Login a user
	loggedUser, err := s.Login(context.TODO(), "srojas@gmail.com")
	log.Printf("%v, %v\n\n\n", loggedUser.User, loggedUser.Token)

	// Verify access token
	id, err := s.VerifyToken(context.TODO(), loggedUser.Token)

	// id, err := s.VerifyToken(context.TODO(), &entities.Token{
	// 	Access:  "a-token",
	// 	Refresh: "r-token",
	// })
	log.Printf("%v -- %v\n", id, err)

	// Refresh access token
	// token, err := s.RefreshToken(context.TODO(), &entities.Token{
	// 	Access:  "a-token",
	// 	Refresh: "r-token",
	// })
	// log.Printf("%v -- %v\n", token, err)

	// Logout
	s.Logout(context.TODO(), &entities.Token{
		Access:  "a-token",
		Refresh: "r-token",
	})
}
