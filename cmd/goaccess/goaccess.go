package main

import (
	"context"
	"os"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/utils"

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
	// loggedUser, err := s.Login(context.TODO(), "srojas@gmail.com")
	// log.Printf("%v, %v\n\n\n", loggedUser.User, loggedUser.Token)

	// Verify access token
	//id, err := s.VerifyToken(context.TODO(), loggedUser.Token)

	// id, err := s.VerifyToken(context.TODO(), &entities.Token{
	// 	Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0aWkzbjNjMW9zczI4aW0zdHZnIiwiZXhwIjoxNjAwNDYzOTI0LCJ1c2VyX2lkIjoiMSJ9.46zww2J_pdAMIDdQKICFwKDhPFwmXnchPgvefCr5LmI",
	// 	Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA0NjQ1MjQsInJlZnJlc2hfdXVpZCI6ImJ0aWkzbjNjMW9zczI4aW0zdTAwIiwidXNlcl9pZCI6IjEifQ.8F1c-Py_tWH_bxa0taF68TftihzL-WhPsUV7WiJcSCA",
	// })
	// log.Printf("%v -- %v\n", id, err)

	// Refresh access token
	// token, err := s.RefreshToken(context.TODO(), &entities.Token{
	// 	Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0aWkzbjNjMW9zczI4aW0zdHZnIiwiZXhwIjoxNjAwNDYzOTI0LCJ1c2VyX2lkIjoiMSJ9.46zww2J_pdAMIDdQKICFwKDhPFwmXnchPgvefCr5LmI",
	// 	Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA0NjQ1MjQsInJlZnJlc2hfdXVpZCI6ImJ0aWkzbjNjMW9zczI4aW0zdTAwIiwidXNlcl9pZCI6IjEifQ.8F1c-Py_tWH_bxa0taF68TftihzL-WhPsUV7WiJcSCA",
	// })
	// log.Printf("%v -- %v\n", token, err)

	// Logout
	s.Logout(context.TODO(), &entities.Token{
		Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0aWk1Z2JjMW9zczZwazBnYjFnIiwiZXhwIjoxNjAwNDY0MTUzLCJ1c2VyX2lkIjoiMSJ9.vNsS85rc4WzUttmQ6YnByGbJFJPlqHf5e8Slmv_1FgM",
		Refresh: "r-token",
	})

	initRepo, err := repository.NewInitRepository(ctx, redisClient)
	path, _ := os.Getwd()
	path = path + "/init"
	jsonHandler := utils.NewJSONHandler(path)
	initService := service.NewInitService(initRepo, jsonHandler)
	initService.Init(ctx, true)
}
