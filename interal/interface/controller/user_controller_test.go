package controller

import (
	"context"
	"net/http"
	"testing"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	mock_interfactor "git.foxminded.com.ua/3_REST_API/interal/usecase/interactor/mock_interfactor"
)

func TestSignUpHandler(t *testing.T) {
	type mockBehavior func(*mock_interfactor.MockUserInteractor, context.Context, *models.User)

	testTable := []struct {
		name      string
		inputBody string
		inputMock struct {
			context.Context
			models.User
		}
		mockBehavior                     mockBehavior
		expectedStatusCode               int
		expectedResponseBody             string
		expectedCookieAuthorizationValue string
		expectedCookieDuration           int
	}{
		{
			name:      "OK",
			inputBody: `"user_name": "username", "first_name": "Bogdan", "last_name": "Pavliuk", "password": "WeryDifficultpassword(1234)"`,
			inputMock: struct {
				context.Context
				models.User
			}{context.Background(), models.User{UserName: "username", FirstName: "Bogdan", LastName: "Pavliuk", Password: "WeryDifficultpassword(1234)"}},
			mockBehavior: func(mui *mock_interfactor.MockUserInteractor, ctx context.Context, u *models.User) {
				mui.EXPECT().SignUp(ctx, u).Return(86400, "token", nil)
			},
			expectedStatusCode:               http.StatusOK,
			expectedResponseBody:             `"token": "token", "message": "You are logged in!"`,
			expectedCookieAuthorizationValue: "token",
			expectedCookieDuration:           86400,
		},
		{
			name:               "Wrong Input",
			inputBody:          `"user_name", "first_name", "last_name"`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: "couldn't bind some data : [code=400, message=Syntax error: offset=13, error=invalid character ',' after object key, " +
				"internal=invalid character ',' after object key]",
			expectedCookieAuthorizationValue: "",
			expectedCookieDuration:           0,
		},
		{
			name:               "Validation error",
			inputBody:          `"user_name": "username", "first_name": "Bogdan", "last_name": "Pavliuk", "password": "1234"`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: "validation cannot be passed : [password must containe at least 7 letters,  1 number, 1 upper case, 1 special character. " +
				" err: Key: 'SignUpRequest.Password' Error:Field validation for 'Password' failed on the 'password' tag]",
			expectedCookieAuthorizationValue: "",
			expectedCookieDuration:           0,
		},
		{
			name:      "exist different user with this username",
			inputBody: `"user_name": "username", "first_name": "Bogdan", "last_name": "Pavliuk", "password": "WeryDifficultpassword(1234)"`,
			inputMock: struct {
				context.Context
				models.User
			}{context.Background(), models.User{UserName: "username", FirstName: "Bogdan", LastName: "Pavliuk", Password: "WeryDifficultpassword(1234)"}},
			mockBehavior: func(mui *mock_interfactor.MockUserInteractor, ctx context.Context, u *models.User) {
				mui.EXPECT().SignUp(ctx, u).Return(86400, "token", nil)
			},
		},
	}
}
