package users

type CreateUserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserExistsDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

const (
	UserNotFound = "user not found"
	UserExists   = "user already exists"
)
