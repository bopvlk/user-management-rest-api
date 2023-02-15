package requests

type SignUpRequest struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type SignInRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
