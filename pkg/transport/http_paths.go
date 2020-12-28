package transport

import (
	"fmt"
)

const accessPrefix = "/api/access"
const authorizationPrefix = "/api/auth"

var userPaths = map[string]string{
	"registerUser":   "/users",              // gRPC
	"unregisterUser": "/users",              // gRPC
	"login":          "/users/login",        // gRPC
	"logout":         "/users/logout",       // gRPC
	"verify":         "users/token/verify",  // gRPC
	"refresh":        "users/token/refresh", // gRPC
}

var accessPaths = map[string]string{
	"addRole":            "/roles",                                              //POST
	"editRole":           "/roles",                                              //PUT
	"deleteRole":         "/roles",                                              //DELETE
	"assignModules":      "/modules/{role_id}",                                  //POST
	"unassignModules":    "/modules/{role_id}/{modules}",                        //DELETE
	"assignSubModules":   "/submodule/{role_id}/modules/{module}",               //POST
	"unassignSubModules": "/submodule/{role_id}/modules/{module}/{sub_modules}", //DELETE
	"assignSections":     "/section/{role_id}/modules/{module}",                 //POST
	"unassignSections":   "/section/{role_id}/modules/{module}/{sections}",      //DELETE
}

var authorizationPaths = map[string]string{
	"assignActions": "/actions", //POST
}

func getAccessPath(path string) string {
	if p, ok := accessPaths[path]; ok {
		fmt.Println(accessPrefix + p)
		return accessPrefix + p
	}
	panic("Undefined path")
}
