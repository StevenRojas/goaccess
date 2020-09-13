package main

import (
	"context"

	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/repository"

	"github.com/StevenRojas/goaccess/pkg/service"

	"github.com/StevenRojas/goaccess/pkg/configuration"
	"github.com/StevenRojas/goaccess/pkg/redis"
)

func main() {
	ctx := context.Background()
	serviceConfig, err := configuration.Read()
	if err != nil {
		panic(err)
	}
	redisClient, err := redis.NewRedisClient(ctx, serviceConfig.Redis)
	if err != nil {
		panic(err)
	}
	usersRepo, err := repository.NewUsersRepository(ctx, redisClient)
	s := service.NewAccessService(usersRepo, serviceConfig.Security)

	// Login a user
	// loggedUser, err := s.Login(context.TODO(), "srojas@gmail.com")
	// log.Printf("%v, %v\n\n\n", loggedUser.User, loggedUser.Token)

	// Verify access token
	// id, err := s.VerifyToken(context.TODO(), loggedUser.Token)
	// id, err := s.VerifyToken(context.TODO(), &entities.Token{
	// 	Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0ZW1hNTNjMW9zamNpdXM1ZG9nIiwiZXhwIjoxNTk5OTU2ODQ0LCJ1c2VyX2lkIjoiMSJ9.vUt1zUk-tBEKMy92ycUNANcb4MQfrfNaJyGvjVj_zts",
	// 	Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTk5NTM3OTgsInJlZnJlc2hfdXVpZCI6ImJ0ZWxkbGpjMW9zaG9tbms1N3UwIiwidXNlcl9pZCI6IjEifQ.v4J5AI91de-ONr6a8tc37cpTYni5wPxVKW5wYZ4qNPU",
	// })
	// log.Printf("%v -- %v\n", id, err)

	// token, err := s.RefreshToken(context.TODO(), &entities.Token{
	// 	Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0ZW1pcmJjMW9za2pzMWZwZHEwIiwiZXhwIjoxNTk5OTU3OTU3LCJ1c2VyX2lkIjoiMSJ9.2yYSFdtZM6uZ_iwAs0NNkIMbbuBFVc4F1yhGKokd3fA",
	// 	Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTk5NTg1NTcsInJlZnJlc2hfdXVpZCI6ImJ0ZW1pcmJjMW9za2pzMWZwZHFnIiwidXNlcl9pZCI6IjEifQ.tIiItjiKnS1QdUaSYj55ObVFNMPtYLtiI8g48GC8azo",
	// })
	// log.Printf("%v -- %v\n", token, err)

	s.Logout(context.TODO(), &entities.Token{
		Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ0ZW1qbGJjMW9za21mMjluYmRnIiwiZXhwIjoxNTk5OTU4MDYxLCJ1c2VyX2lkIjoiMSJ9.sS2VvzfNZiicgEzDyIFWsmrX4QTdydlLd9G63GnvQRY",
		Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTk5NTg2NjEsInJlZnJlc2hfdXVpZCI6ImJ0ZW1qbGJjMW9za21mMjluYmUwIiwidXNlcl9pZCI6IjEifQ.ojAGrS4T3EIUJqxZSuATI87TQ5ghHthQL3n0uuq0IDg",
	})
}
