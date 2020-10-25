package repository

const configKey string = "config:"
const accessTemplateKey string = "config:access"
const actionTemplateKey string = "config:action"
const isSetKey string = "isset"

const userKey string = "user:%s"         // user:userID
const usersKey string = "users"          // users
const roleUserKey string = "roleuser:%s" // roleuser:roleID
const userRoleKey string = "userrole:%s" // userrole:userID

const roleIDKey string = "roleId"
const rolesKey string = "roles"

const actionsModuleKey string = "%s:%s:mo"        // rolesKey:roleID:mo
const roleActionsKey string = "%s:%s:ac:%s:%s"    // rolesKey:roleID:ac:moduleName:submoduleName
const actionsByModuleKey string = "actions:%s:%s" // actions:userID:moduleName
const hasPesmissionKey string = "actionlist:%s"   // actionlist:actionName

const roleModulesKey string = "%s:%s:mo"        // rolesKey:roleID:mo
const roleSubModulesKey string = "%s:%s:sm:%s"  // rolesKey:roleID:sm:moduleName
const roleSectionsKey string = "%s:%s:se:%s:%s" // rolesKey:roleID:se:moduleName:submoduleName
const accessKey string = "access:%s"            // access:userID
