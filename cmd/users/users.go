package users

type User struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
