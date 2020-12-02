package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/StevenRojas/goaccess/pkg/events"
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
	modulesRepo, err := repository.NewModulesRepository(ctx, redisClient)
	if err != nil {
		panic(err)
	}
	rolesRepo, err := repository.NewRolesRepository(ctx, redisClient)
	if err != nil {
		panic(err)
	}
	actionsRepo, err := repository.NewActionsRepository(ctx, redisClient)
	if err != nil {
		panic(err)
	}
	initRepo, err := repository.NewInitRepository(ctx, redisClient)
	if err != nil {
		panic(err)
	}

	jwtHander := utils.NewJwtHandler(serviceConfig.Security)
	s := service.NewAuthenticationService(usersRepo, jwtHander)

	// Register a user
	// err = s.Register(context.TODO(), &entities.User{
	// 	ID:      "1",
	// 	Email:   "steven.rojas@gmail.com",
	// 	Name:    "Steven Rojas",
	// 	IsAdmin: false,
	// })
	// log.Printf("%v\n", err)

	// Unregister a user
	// err = s.Unregister(context.TODO(), &entities.User{
	// 	ID:      "1",
	// 	Email:   "steven.rojas@gmail.com",
	// 	Name:    "Steven Rojas",
	// 	IsAdmin: false,
	// })
	// log.Printf("%v\n", err)

	// Login a user
	// loggedUser, err := s.Login(context.TODO(), "steven.rojas@gmail.com")
	// log.Printf("%v, %v\n\n\n", loggedUser.User, loggedUser.Token)

	// Verify access token
	// id, err := s.VerifyToken(context.TODO(), loggedUser.Token.Access)

	id, _ := s.VerifyToken(
		context.TODO(),
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ1YTdrZ2JjMW9zZ3JiYTFlMTYwIiwiZXhwIjoxNjAzNTY2MjY1LCJ1c2VyX2lkIjoiMSJ9.HvDINzWKdnqg5xBNhXh8e8jlVZfG-KTuL_cqN_5mCQ0",
	)
	log.Printf("%v\n", id)

	// Refresh access token
	// token, err := s.RefreshToken(
	// 	context.TODO(),
	// 	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM1NzUwNTksInJlZnJlc2hfdXVpZCI6ImJ1YTltcnJjMW9zajIxOGZyZTRnIiwidXNlcl9pZCI6IjEifQ.YSXm392ugCnQbz2VCn-bdGQZaINxgxi0E6-9EYjqr5Y",
	// )
	// log.Printf("%v -- %v\n", token, err)

	// Logout
	// s.Logout(context.TODO(), &entities.Token{
	// 	Access:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJ1YTdrZ2JjMW9zZ3JiYTFlMTYwIiwiZXhwIjoxNjAzNTY2MjY1LCJ1c2VyX2lkIjoiMSJ9.HvDINzWKdnqg5xBNhXh8e8jlVZfG-KTuL_cqN_5mCQ0",
	// 	Refresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDM1NjY1NjUsInJlZnJlc2hfdXVpZCI6ImJ1YTdrZ2JjMW9zZ3JiYTFlMTZnIiwidXNlcl9pZCI6IjEifQ.5NT3S35Y8Df43LDNxKydPkZRIQnd-179h4fcBtSytFg",
	// })

	path, _ := os.Getwd()
	path = path + "/init"
	jsonHandler := utils.NewJSONHandler(path)
	initService := service.NewInitService(initRepo, jsonHandler)
	initService.Init(ctx, false)

	// Role events listener
	subscriberFeed := events.NewSubscriber()
	accessListener := events.NewAccessListener(modulesRepo, rolesRepo, subscriberFeed)
	go accessListener.RegisterAccessListener()
	actionListener := events.NewActionListener(actionsRepo, rolesRepo, subscriberFeed)
	go actionListener.RegisterActionListener()

	// Handle access
	//accessService := service.NewAccessService(modulesRepo, rolesRepo, actionsRepo, subscriberFeed)

	// Add role
	// roleID, err := accessService.AddRole(ctx, "accounting manager")
	// log.Printf("%v -- %v\n", roleID, err)

	// Update role
	// err = accessService.EditRole(ctx, "r1", "accounting manager updated")
	// log.Printf("%v\n", err)

	// Delete role
	// err = accessService.DeleteRole(ctx, "r2")
	// log.Printf("%v\n", err)

	// Assign modules
	// err = accessService.AssignModules(ctx, "r1", []string{"bank"})
	// log.Printf("%v\n", err)
	// Unassign modules
	//err = accessService.UnassignModules(ctx, "r1", []string{"bank"})
	//log.Printf("%v\n", err)

	// Assign submodules
	// err = accessService.AssignSubModules(ctx, "r1", "bank", []string{"accounts"})
	//log.Printf("%v\n", err)

	// Unassign submodules
	// err = accessService.UnassignSubModules(ctx, "r1", "vehicles", []string{"work-subcategory"})
	// log.Printf("%v\n", err)

	// Assign sections
	// err = accessService.AssignSections(ctx, "r1", "vehicles", "reception", []string{"finder", "add", "test"})
	// log.Printf("%v\n", err)

	// Unassign sections
	// err = accessService.UnassignSections(ctx, "r1", "vehicles", "reception", []string{"test"})
	// log.Printf("%v\n", err)

	//moduleList, err := accessService.ModulesList(ctx)
	//log.Printf("%v\n", moduleList)

	// Get module structure for a new role
	// module, err := accessService.ModuleStructure(ctx, "vehicles")
	// j, _ := json.Marshal(module)
	// fmt.Printf("%v\n", string(j))

	// assignations, err := accessService.GetRoleAccessList(ctx, "r1")
	// log.Printf("error %v\n", err)
	// fmt.Printf("%v\n", assignations)

	// AuthorizationService
	authorizationService := service.NewAuthorizationService(modulesRepo, rolesRepo, actionsRepo, usersRepo, subscriberFeed)

	// Actions
	// err = authorizationService.AssignActions(ctx, "r1", "bank", "accounts", []string{"delete:account:[]"})
	// log.Printf("%v\n", err)
	// err = authorizationService.UnassignActions(ctx, "r1", "vehicles", "reception", []string{"delete:brand:[]:test"})
	// log.Printf("%v\n", err)

	// Assign roles to users
	//err = authorizationService.AssignRole(ctx, "1", "r1")
	// log.Printf("%v\n", err)
	//err = authorizationService.UnassignRole(ctx, "1", "r1")
	//log.Printf("%v\n", err)

	// accessJSON, err := authorizationService.GetAccessList(ctx, "1")
	// log.Printf("%v\n", err)
	// log.Printf("%v\n", accessJSON)

	//actionsJSON, err := authorizationService.GetActionListByModule(ctx, "bank", "1")
	// log.Printf("%v\n", err)
	//log.Printf("%v\n", actionsJSON)
	hasPermission, err := authorizationService.CheckPermission(ctx, "delete:account:[]", "1")
	log.Printf("%v\n", err)
	log.Printf("%v\n", hasPermission)

	// fmt.Println("*************************\n\n\n")
	// //path = path + "/postman"
	// path = "/Users/steven.rojas/code/go/goaccess/init/postman"
	// postman := utils.NewPostmanParser(path, false)
	// j, err := postman.Parse("revicart_collection.json", "vehicles")
	// log.Printf("%v\n", err)
	// log.Printf("%v\n", j)

	time.Sleep(3 * time.Second)
}
