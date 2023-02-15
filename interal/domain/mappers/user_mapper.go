package mappers

import (
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/domain/requests"
)

func MapUsersModelToUsersResponse(u []*models.User) []*requests.UserResponse {
	ur := make([]*requests.UserResponse, len(u))
	for i := 0; i < len(u); i++ {
		ur[i] = MapUserModelToUserResponse(u[i])
	}
	return ur
}

func MapUserModelToUserResponse(u *models.User) *requests.UserResponse {
	return &requests.UserResponse{
		ID:        u.ID,
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: &u.DeletedAt.Time,
	}
}

func MapSignUpRequestToUserModel(signUp *requests.SignUpRequest) *models.User {
	return &models.User{
		UserName:  signUp.UserName,
		FirstName: signUp.FirstName,
		LastName:  signUp.LastName,
		Password:  signUp.Password,
	}

}

func MapSignInRequestToUserModel(signIp *requests.SignInRequest) *models.User {
	return &models.User{
		UserName: signIp.UserName,
		Password: signIp.Password,
	}

}
