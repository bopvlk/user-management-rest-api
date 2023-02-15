package requests

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
)

type SignUpInResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	IsError bool   `json:"is_error"`
}

type GetUsersResponse struct {
	Message       string             `json:"message"`
	UsersResponse *models.Pagination `json:"users"`
	IsError       bool               `json:"is_error"`
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
