package interactor

import (
	"context"
	"git.foxminded.com.ua/3_REST_API/gen/mocks"
	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestFindOneSigner(t *testing.T) {
	now := time.Now()
	testTable := []struct {
		scenario      string
		expectedUser  *models.User
		expectedError error
	}{
		{
			"find one user by id",
			&models.User{
				ID:        121,
				UserName:  "JohnHall",
				FirstName: "John",
				LastName:  "Hall",
				Password:  "1231",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			nil,
		},
		{
			"user not found by id",
			&models.User{},
			&apperrors.UserNotFoundErr,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := &userInteractor{
		userRepo:       userRepoMock,
		hashSalt:       "",
		signingKey:     nil,
		expireDuration: 0,
	}

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.Background()
			userRepoMock.EXPECT().FindOneUserByID(ctx, tc.expectedUser.ID).Return(tc.expectedUser, tc.expectedError)
			user, err := uInteractor.FindOneSigner(ctx, tc.expectedUser.ID)
			if err != nil {

				if tc.expectedError != nil && apperrors.Is(err, tc.expectedError.(*apperrors.AppError)) {
					return
				}

				t.Fatal(err)
			}

			assert.Equal(t, user, tc.expectedUser)
		})
	}
}
