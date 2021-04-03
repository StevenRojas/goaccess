package codec

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/StevenRojas/goaccess/pkg/entities"
	e "github.com/StevenRojas/goaccess/pkg/errors"
	"github.com/gorilla/mux"
)

// DecodeAddRoleRequest decode request
func DecodeAddRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	roleRequest := entities.Role{}

	errs := entities.InitRoleValidator(r, &roleRequest).ValidateJSON()
	if len(errs) > 0 {
		return nil, e.ValidationError{
			Message: "Validation error while adding a Role",
			Fields:  errs,
		}
	}
	return roleRequest, nil
}

// DecodeEditRoleRequest decode request
func DecodeEditRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return nil, e.HTTPBadRequestFromString("Role ID is missing")
	}
	defer r.Body.Close()
	roleRequest := entities.Role{}

	errs := entities.InitRoleValidator(r, &roleRequest).ValidateJSON()
	if len(errs) > 0 {
		return nil, e.ValidationError{
			Message: "Validation error while editing a Role",
			Fields:  errs,
		}
	}
	roleRequest.ID = roleID
	return roleRequest, nil
}

// DecodeDeleteRoleRequest decode request
func DecodeDeleteRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return nil, e.HTTPBadRequestFromString("Role ID is missing")
	}
	return roleID, nil
}

// DecodeGetRoleRequest decode request
func DecodeGetRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	return roleID, nil
}

// DecodeGetRolesByUserRequest decode request
func DecodeGetRolesByUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("User ID is missing")
	}
	return userID, nil
}

// DecodeAssignModulesRequest decode request
func DecodeAssignModulesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	defer r.Body.Close()
	moduleList := entities.ModuleList{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&moduleList)
	if err != nil {
		return nil, e.HTTPBadRequestFromString("Validation error while assigning modules")
	}
	moduleList.RoleID = roleID
	return moduleList, nil
}

// DecodeUnassignModulesRequest decode request
func DecodeUnassignModulesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	modules, ok := vars["modules"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Modules are missing")
	}
	defer r.Body.Close()
	moduleList := entities.ModuleList{}
	moduleList.RoleID = roleID
	moduleList.Modules = strings.Split(modules, ",")
	for i := range moduleList.Modules {
		moduleList.Modules[i] = strings.TrimSpace(moduleList.Modules[i])
	}
	return moduleList, nil
}

// DecodeAssignSubModulesRequest decode request
func DecodeAssignSubModulesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	defer r.Body.Close()
	submoduleList := entities.SubModuleList{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&submoduleList)
	if err != nil {
		return nil, e.HTTPBadRequestFromString("Validation error while assigning submodules")
	}
	submoduleList.RoleID = roleID
	submoduleList.Module = moduleName
	return submoduleList, nil
}

// DecodeUnassignSubModulesRequest decode request
func DecodeUnassignSubModulesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	submodules, ok := vars["sub_modules"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Submodules are missing")
	}
	defer r.Body.Close()
	submoduleList := entities.SubModuleList{}
	submoduleList.RoleID = roleID
	submoduleList.Module = moduleName
	submoduleList.SubModules = strings.Split(submodules, ",")
	for i := range submoduleList.SubModules {
		submoduleList.SubModules[i] = strings.TrimSpace(submoduleList.SubModules[i])
	}
	return submoduleList, nil
}

// DecodeAssignSectionsRequest decode request
func DecodeAssignSectionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	submoduleName, ok := vars["sub_module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Submodule name is missing")
	}
	defer r.Body.Close()
	sectionList := entities.SectionList{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&sectionList)
	if err != nil {
		return nil, e.HTTPBadRequestFromString("Validation error while assigning sections")
	}
	sectionList.RoleID = roleID
	sectionList.Module = moduleName
	sectionList.SubModule = submoduleName
	return sectionList, nil
}

// DecodeUnassignSectionsRequest decode request
func DecodeUnassignSectionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	submoduleName, ok := vars["sub_module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Subodule name is missing")
	}
	sections, ok := vars["sections"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Sections are missing")
	}
	sectionList := entities.SectionList{}
	sectionList.RoleID = roleID
	sectionList.Module = moduleName
	sectionList.SubModule = submoduleName
	sectionList.Sections = strings.Split(sections, ",")
	for i := range sectionList.Sections {
		sectionList.Sections[i] = strings.TrimSpace(sectionList.Sections[i])
	}
	return sectionList, nil
}
