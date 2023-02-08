package mappers

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
)

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

type SignUpInResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	IsError bool   `json:"is_error"`
}

type GetUsersResponse struct {
	Message       string          `json:"message"`
	UsersResponse []*UserResponse `json:"users"`
	IsError       bool            `json:"is_error"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	UserName  string     `json:"user_name"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GetOneUserResponse struct {
	Message      string       `json:"message"`
	UserResponse UserResponse `json:"user"`
	IsError      bool         `json:"is_error"`
}

func MapUsersModelToUsersResponse(u []*models.User) []*UserResponse {
	ur := make([]*UserResponse, len(u))
	for i := 0; i < len(u); i++ {
		// ur[i].ID = u[i].ID
		// ur[i].UserName = u[i].UserName
		// ur[i].FirstName = u[i].FirstName
		// ur[i].LastName = u[i].LastName
		// ur[i].CreatedAt = u[i].CreatedAt
		// ur[i].DeletedAt = &u[i].DeletedAt.Time
		ur[i] = MapUserModelToUserResponse(u[i])
	}
	return ur
}

func MapUserModelToUserResponse(u *models.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: &u.DeletedAt.Time,
	}
}

func MapSignUpRequestToUserModel(signUp *SignUpRequest) *models.User {
	return &models.User{
		UserName:  signUp.UserName,
		FirstName: signUp.FirstName,
		LastName:  signUp.LastName,
		Password:  signUp.Password,
	}

}

func MapSignInRequestToUserModel(signIp *SignInRequest) *models.User {
	return &models.User{
		UserName: signIp.UserName,
		Password: signIp.Password,
	}

}
