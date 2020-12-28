package entities

const (
	EventTypeAccess = "EventTypeAccess"
	EventTypeAction = "EventTypeAction"
)

// User struct
type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

// Token struct
type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// LoggedUser logged user struct
type LoggedUser struct {
	User  *User
	Token *Token
}

type SubModule struct {
	Name     string            `json:"submodule"`
	Access   bool              `json:"access"`
	Actions  map[string]Action `json:"actions,omitempty"`
	Sections map[string]bool   `json:"sections,omitempty"`
}

type Module struct {
	Name       string      `json:"module"`
	Access     bool        `json:"access"`
	SubModules []SubModule `json:"submodules"`
}

type SubModuleInit struct {
	Name        string            `json:"submodule"`
	SectionList []string          `json:"sectionList,omitempty"`
	ActionList  map[string]string `json:"actionList,omitempty"`
	Actions     map[string]Action `json:"actions"`
	Sections    map[string]bool   `json:"sections"`
}

type ModuleInit struct {
	Name       string          `json:"module"`
	SubModules []SubModuleInit `json:"submodules"`
}

type Action struct {
	Title   string `json:"title"`
	Allowed bool   `json:"allowed"`
}

type ActionSubModule struct {
	Name       string            `json:"submodule"`
	ActionList map[string]string `json:"actionList,omitempty"`
	Actions    map[string]Action `json:"actions"`
}

type ActionModule struct {
	Name       string            `json:"module"`
	SubModules []ActionSubModule `json:"submodules"`
}

type RoleEvent struct {
	RoleID    string
	UserID    string
	EventType string
}
