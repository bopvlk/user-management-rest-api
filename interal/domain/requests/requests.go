package requests

type SignUpRequest struct {
	UserName  string `json:"user_name" validate:"required,min_len=5"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,password"`
}

type SignInRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
