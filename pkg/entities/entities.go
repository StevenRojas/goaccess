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
