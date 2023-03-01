package requests

type SignUpRequest struct {
	UserName  string `json:"user_name" validate:"required,min=5"`
	Role      string `json:"role" validate:"required,contains=user|contains=moderator|contains=admin"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,password,min=7"`
}

type SignInRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	UserName  string `json:"user_name" validate:"required,min=5"`
	Role      string `json:"role" validate:"required,contains=user|moderator|admin"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
