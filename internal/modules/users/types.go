package users

type CreateUserDto struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserExistsDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

const (
	UserNotFound = "user not found"
	UserExists   = "user already exists"
)
