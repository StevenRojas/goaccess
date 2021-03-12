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
	"addRole":               "/roles",           //POST
	"editRole":              "/roles/{role_id}", //PUT
	"deleteRole":            "/roles/{role_id}",
	"getAllModules":         "/modules",                                                                //GET
	"getAssignedModules":    "/modules/{role_id}",                                                      //GET
	"assignModules":         "/modules/{role_id}",                                                      //POST
	"unassignModules":       "/modules/{role_id}/{modules}",                                            //DELETE
	"getAssignedSubModules": "/submodules/{role_id}",                                                   //GET
	"assignSubModules":      "/submodules/{role_id}/modules/{module}",                                  //POST
	"unassignSubModules":    "/submodules/{role_id}/modules/{module}/{sub_modules}",                    //DELETE
	"getAssignedSections":   "/sections/{role_id}",                                                     //GET
	"assignSections":        "/sections/{role_id}/modules/{module}/submodules/{sub_module}",            //POST
	"unassignSections":      "/sections/{role_id}/modules/{module}/submodules/{sub_module}/{sections}", //DELETE
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
