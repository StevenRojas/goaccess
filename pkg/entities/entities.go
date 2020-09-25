package entities

// User struct
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
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
	Name     string   `json:"submodule"`
	Sections []string `json:"sections"`
}

type Module struct {
	Name       string      `json:"module"`
	SubModules []SubModule `json:"submodules"`
}

type ActionSubModule struct {
	Name    string   `json:"submodule"`
	Actions []string `json:"actions"`
}

type ActionModule struct {
	Name       string            `json:"module"`
	SubModules []ActionSubModule `json:"submodules"`
}
