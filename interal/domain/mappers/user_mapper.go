package mappers

import (
	"fmt"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/domain/requests"
	"github.com/labstack/echo/v4"
)

func MapUserToUserResponse(u *models.User) *requests.UserResponse {
	return &requests.UserResponse{
		ID:        u.ID,
		UserName:  u.UserName,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: &u.DeletedAt.Time,
	}
}

func MapSignUpRequestToUser(signUp *requests.SignUpRequest) *models.User {
	return &models.User{
		UserName:  signUp.UserName,
		Role:      signUp.Role,
		FirstName: signUp.FirstName,
		LastName:  signUp.LastName,
		Password:  signUp.Password,
	}
}

func MapUpdateRequestToUser(signUp *requests.UpdateRequest) *models.User {
	return &models.User{
		UserName:  signUp.UserName,
		Role:      signUp.Role,
		FirstName: signUp.FirstName,
		LastName:  signUp.LastName,
	}
}

func MapUpdateOwnRequestToUser(signUp *requests.UpdateOwnRequest) *models.User {
	return &models.User{
		UserName:  signUp.UserName,
		Role:      signUp.Role,
		FirstName: signUp.FirstName,
		LastName:  signUp.LastName,
		Password:  signUp.Password,
	}
}

func MapSignInRequestToUser(signIp *requests.SignInRequest) *models.User {
	return &models.User{
		UserName: signIp.UserName,
		Password: signIp.Password,
	}

}

func MapAppErrorToHTTPError(err error) *echo.HTTPError {
	appErr := err.(*apperrors.AppError)
	return echo.NewHTTPError(appErr.HTTPCode, appErr.Error())
}

func MapUserToGetUserResponse(user *models.User) *requests.GetOneUserResponse {
	return &requests.GetOneUserResponse{
		Message:      fmt.Sprintf("There is user with ID %v", user.ID),
		UserResponse: MapUserToUserResponse(user),
	}
}

func MapPaginationAndUsersToGetUsersResponse(users []*models.User, pagination *models.Pagination, name string) *requests.GetUsersResponse {
	message := fmt.Sprintf("Hello,%v U are in restricted zone.", name)
	ur := make([]*requests.UserResponse, len(users))
	for i := 0; i < len(users); i++ {
		ur[i] = MapUserToUserResponse(users[i])

	}

	pagination.Rows = ur
	return &requests.GetUsersResponse{
		Message:       message,
		UsersResponse: pagination,
	}
}

func MapUserToUpdateResponse(u *models.User) *requests.UpdateUserResponce {
	return &requests.UpdateUserResponce{
		Message: fmt.Sprintf("There is updated user with id: %d", u.ID),
		UserResponse: requests.UserResponse{
			ID:        u.ID,
			UserName:  u.UserName,
			Role:      u.Role,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: &u.DeletedAt.Time,
		},
	}
}
