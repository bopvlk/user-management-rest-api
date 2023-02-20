package controller

import (
	"testing"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	mock_interfactor "git.foxminded.com.ua/3_REST_API/interal/usecase/interactor/mock_interfactor"
)

func TestSignUpHandler(t *testing.T) {
	type mockBehavior func(*mock_interfactor.MockUserInteractor, *models.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            *models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{name: "OK",
			inputBody: `"user_name": "username", "first_name": "Bogdan", "last_name": "Pavliuk", "password": "werydifficultpassword(1234)"`},
	}
}
