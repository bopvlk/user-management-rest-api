package requests

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
)

type SignUpInResponse struct {
	Message string `json:"message"`
}

type GetUsersResponse struct {
	Message       string             `json:"message"`
	UsersResponse *models.Pagination `json:"users"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	UserName  string     `json:"user_name"`
	Role      string     `json:"role"`
	Rating    int        `json:"rating"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GetOneUserResponse struct {
	Message      string        `json:"message"`
	UserResponse *UserResponse `json:"user"`
}

type UpdateUserResponce struct {
	Message      string        `json:"message"`
	UserResponse *UserResponse `json:"user"`
}

type ErrorResponse struct {
	Message  string `json:"message"`
	HTTPCode int    `json:"err_code"`
}
