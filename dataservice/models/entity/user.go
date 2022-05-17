package entity

type User struct {
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserPassword string `json:"user_password"`
	UserSalt     string `json:"user_salt"`
}
