

# goaccess

Go package to grant and check access to users with JWT, roles, access and permissions. The principle of `goaccess` package is that a web site or app is structured based on `modules`, `submodules` and `sections`. For example the `modules` could be the options in the main menu, the `submodules` are the options for each menu and the `sections` are the rendered page sections. `goaccess` handle the structure in a JSON with this schema:
```json
{
  "vehicles": {
    "module": "vehicles",
    "access": true,
    "submodules": [
      {
        "submodule": "brand",
        "access": false,
        "sections": {
          
        }
      },
      {
        "submodule": "reception",
        "access": true,
        "sections": {
          "add": false,
          "finder": true
        }
      }
    ]
  }
}
```

On the other hand, `goaccess` package has the concept of `actions` in order to handle user permissions. `actions` are related to the `submodules` and are plain strings (whatever string that make sense for identify permissions) that could be active or not. For example, you can use strings based on the API endpoints the web app will use, here an schema example:
```json
{
  "module": "vehicles",
  "submodules": [
    {
      "submodule": "brand",
      "actionList": {
        "delete:brand:[]": "Delete brand",
        "delete:brand:[]:remove": "Delete brand remove",
        "patch:brand:[]:restore": "Modify brand restore",
        "post:brand": "Create brand",
        "put:brand:[]": "Update brand"
      }
    },
    {
      "submodule": "vehicle",
      "actionList": {
        "delete:vehicle:[]:photo:[]": "Delete vehicle photo",
        "post:vehicle": "Create vehicle",
        "post:vehicle:[]:photo": "Create vehicle photo",
        "put:vehicle:[]": "Update vehicle"
      }
    }
  ]
}
```
where `:` replace the `/` and `[]` is a placeholder for the URL params. It is not necessary to have `get` endpoints in the list, because they are controlled on the `module > submodule > section` access.

## Installation
```go
go get github.com/StevenRojas/goaccess
```

## Setup environment
### JWT
```go
export JWT_SECRET_KEY=secret!
export JWT_EXPIRE_HOURS=2
export JWT_REFRESH_HOURS=7
```
### Redis
```go
export REDIS_ADDR=localhost:6379
export REDIS_DB=10
export REDIS_PASS=secret
```

## Initialize Services
Read configuration from environment variables:
```go
serviceConfig, err := configuration.Read()
```
Redis client:
```go
redisClient := redis.NewClient(&redis.Options{
	Addr: serviceConfig.Redis.Addr,
	Password: serviceConfig.Redis.Pass,
	DB: serviceConfig.Redis.DB,
})
```
Repositories:
```go
initRepo, err := repository.NewInitRepository(ctx, redisClient)
usersRepo, err := repository.NewUsersRepository(ctx, redisClient)
modulesRepo, err := repository.NewModulesRepository(ctx, redisClient)
rolesRepo, err := repository.NewRolesRepository(ctx, redisClient)
actionsRepo, err := repository.NewActionsRepository(ctx, redisClient)
```
JWT handler:
```go
jwtHander := utils.NewJwtHandler(serviceConfig.Security)
```
Services (the use of each one is explined at the corresponding sections):
```go
service.NewAuthenticationService(usersRepo, jwtHander)
service.NewInitService(initRepo, jsonHandler)
service.NewAccessService(modulesRepo, rolesRepo, actionsRepo, subscriberFeed)
service.NewAuthorizationService(modulesRepo, rolesRepo, actionsRepo)
```
## Initialization Service

In order to initialize Redis database with modules, submodules, sections and actions is necessary to have a set of json files in folders `init/modules` and `init/actions`, check the examples. Then run the following:
```go
path, _ := os.Getwd()
path = path + "/init"
jsonHandler := utils.NewJSONHandler(path)
initService := service.NewInitService(initRepo, jsonHandler)
initService.Init(ctx, true) // true to force recreate the DB
```
The code above will parse all `json` files at `init` folder and setup Redis.

In order to generate the JSON files, you can use the following utility code that parse a Postman collection and generate the JSON file for a given module:
```go
path, _ := os.Getwd()
path = path + "/postman"
postman := utils.NewPostmanParser(path, true) // use postman description
j, err := postman.Parse("collection.json", "vehicles")
```
where `j` is the JSON schema, `collection.json` is the postman filename and `vehicles` is the module name. The utility will generate a description based on the URL or will use the Postman description attribute. Of course, the JSON generated for this utility should be reviewed and improved, the `sectionList` is empty, so it should be filled out manually.

## Authentication Service
This service handle the logic to handle user authentication using JWT `access` and `refresh` tokens. It is possible to login, logout, vefiry the access token and refresh the token when it expires. 
```go
s := service.NewAuthenticationService(usersRepo, jwtHander)
```

**Register a user:** add a user in the DB, The ID is not autogerated because it suppose there is another module like HR that has a CRUD for users
```go
err = s.Register(context.TODO(), &entities.User{
	ID: "1",
	Email: "steven.rojas@gmail.com",
	Name: "Steven Rojas",
	IsAdmin: false,
})
```
**Unregister a user:** remove user from the DB
```go
err = s.Unregister(context.TODO(), &entities.User{
	ID: "1",
	Email: "steven.rojas@gmail.com",
	Name: "Steven Rojas",
	IsAdmin: false,
})
```
**Login a user:** validates if the user email is active and returns a `user` structure (`id`, `email`, `name`, `admin`) and the `access` and `refresh` tokens. The calims of the `access` token contains only the `user_id`
```json
{
  "access_uuid": "bua7kgbc1osgrba1e160",
  "exp": 1603566265,
  "user_id": "1"
}
```
```go
loggedUser, err := s.Login(context.TODO(), "steven.rojas@gmail.com")
// loggedUser.User
// loggedUser.Token.Access
// loggedUser.Token.Refresh
```
**Verify access token:** Verify if the token is valid and it doesn't expired, returns the `user_id`
```go
id, err := s.VerifyToken(context.TODO(), "c3NfdXVpZ...")
```
**Refresh access token:** Refresh the `access` token after validate that the `refresh` token is still valid, the method returns a new token pair
```go
token, err := s.RefreshToken(context.TODO(), "c3NfdXVpZ...")
```
**Logout:** Logout the user for the given token.
```go
err := s.Logout(context.TODO(), &entities.Token{
	Access: "c3NfdXVpZ...",
	Refresh: "eyJleHAiO...",
})
```
## Access Service
Access service generates events when roles changes, for example when a `module` or an `action` is assigned/unassigned to/from a role. So it is necessary to define subscribers that will update the user's access and permissions as follow:
```go
subscriberFeed := events.NewSubscriber()
accessListener := events.NewAccessListener(subscriberFeed)
go accessListener.RegisterAccessListener() // Listen for access changes
actionListener := events.NewActionListener(subscriberFeed)
go actionListener.RegisterActionListener() // Listen for action changes
```
Then, with the `subscriberFeed` creates the access service instance:
```go
s := service.NewAccessService(modulesRepo, rolesRepo, actionsRepo, subscriberFeed)
```
### Handle roles
With the access service you can add, update and remove roles
```go
// Add role
roleID, err := s.AddRole(ctx, "accounting manager") // roleID = 'r1'
// Update role
err := s.EditRole(ctx, "r1", "accounting manager updated")
// Delete role
err := accessService.DeleteRole(ctx, "r2")
```
In order to get the JSON schema for render modules and actions while creating a new role, use these two methods:
```go
modules, err := s.ModulesForNewRole(ctx)
actions, err := s.ActionsForNewRole(ctx)
```
### Handle modules
You can assign and unassign `modules`, `submodules` and `sections` to a role with the following methods:
```go
// Assign modules
err := s.AssignModules(ctx, "r1", []string{"vehicles"})
// Unassign modules
err := s.UnassignModules(ctx, "r1", []string{"hr", "bank"})
// Assign submodules to a module
err := s.AssignSubModules(ctx, "r1", "vehicles", []string{"reception", "work-category"})
// Unassign submodules
err := s.UnassignSubModules(ctx, "r1", "vehicles", []string{"work-category"})
// Assign sections to a module > submodule
err := s.AssignSections(ctx, "r1", "vehicles", "reception", []string{"finder", "add"})
// Unassign sections
err := s.UnassignSections(ctx, "r4", "vehicles", "reception", []string{"finder"})
```

## Authorization Service
This service allows to assign and unassign roles to/from users, handle role actions and check if a user has permissions to perform an specific action as follow:
```go
s := service.NewAuthorizationService(modulesRepo, rolesRepo, actionsRepo)
// Assign actions to a module > submodule
err := s.AssignActions(ctx, "r1", "vehicles", "brand", []string{"delete:brand:[]:remove"})
// Unassign actions
err := s.UnassignActions(ctx, "r1", "vehicles", "brand", []string{"delete:brand:[]:remove"})
// Assign role to a user
err := s.AssignRole(ctx, "1", "r2")
// Unassign role from a user
err := s.UnassignRole(ctx, "1", "r2")
// Get access JSON for a given user
accessJSON, err := s.GetAccessList(ctx, "1")
// Get action JSON for a given user and module
actionsJSON, err := s.GetActionListByModule(ctx, "vehicle", "1")
// Check if a user has permission to execute an action
hasPermission, err := s.CheckPermission(ctx, "delete|vehicle|brand|[]", "1")
```