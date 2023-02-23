package interactor

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.foxminded.com.ua/3_REST_API/gen/mocks"
	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/interface/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestSignUp(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	testTable := []struct {
		scenario      string
		expectedUser  *models.User
		expectedError error
	}{
		{
			"sing up user",
			&models.User{
				ID:        121,
				UserName:  "JohnHall",
				FirstName: "John",
				LastName:  "Hall",
				Password:  "1234",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			nil,
		},
		{
			"cannot create a user",
			&models.User{
				Password: "1234",
			},
			&apperrors.CanNotCreateUserErr,
		},
		{
			"empty password field",
			&models.User{
				Password: "",
			},
			&apperrors.HashingPasswordErr,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := &userInteractor{
		userRepo:       userRepoMock,
		hashSalt:       "hash_salt",
		signingKey:     []byte("signing_key"),
		expireDuration: 1,
	}

	for _, testCase := range testTable {
		t.Run(testCase.scenario, func(t *testing.T) {
			ctx := context.Background()
			hashingPassword, err := uInteractor.hashing(testCase.expectedUser.Password)
			testCase.expectedUser.Password = hashingPassword
			if err == nil {
				userRepoMock.EXPECT().CreateUser(ctx, testCase.expectedUser).Return(testCase.expectedUser, testCase.expectedError)
			}

			_, tokenString, err := uInteractor.SignUp(ctx, testCase.expectedUser)
			if err != nil {

				if testCase.expectedError != nil && apperrors.Is(err, testCase.expectedError.(*apperrors.AppError)) {
					return
				}
				t.Fatal(err)
			}

			jwtToken, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				return uInteractor.signingKey, nil
			})
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, jwtToken.Claims.(*AuthClaims).User, testCase.expectedUser)
		})
	}
}

func TestSignIn(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	testTable := []struct {
		scenario      string
		inputUserName string
		inputPassword string
		expectedUser  *models.User
		expectedError error
	}{
		{
			"sing ip user",
			"JohnHall",
			"1234",
			&models.User{
				ID:        121,
				UserName:  "JohnHall",
				FirstName: "John",
				LastName:  "Hall",
				Password:  "1234",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			nil,
		},
		{
			"user don't present in database",
			"JohnHall",
			"1234",
			&models.User{},
			&apperrors.UserNotFoundErr,
		},
		{
			"empty password field",
			"JohnHall",
			"",
			&models.User{},
			&apperrors.HashingPasswordErr,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := &userInteractor{
		userRepo:       userRepoMock,
		hashSalt:       "hash_salt",
		signingKey:     []byte("signing_key"),
		expireDuration: 1,
	}

	for _, testCase := range testTable {
		t.Run(testCase.scenario, func(t *testing.T) {
			ctx := context.Background()
			hashingPassword, err := uInteractor.hashing(testCase.inputPassword)
			if err == nil {
				userRepoMock.EXPECT().FindOneUserByUserNameAndPassword(ctx, testCase.inputUserName, hashingPassword).Return(testCase.expectedUser, testCase.expectedError)
			}
			_, token, err := uInteractor.SignIn(ctx, testCase.inputUserName, testCase.inputPassword)
			if err != nil {

				if testCase.expectedError != nil && apperrors.Is(err, testCase.expectedError.(*apperrors.AppError)) {
					return
				}

				t.Fatal(err)
			}

			jwtToken, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				return uInteractor.signingKey, nil
			})
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, jwtToken.Claims.(*AuthClaims).User, testCase.expectedUser)
		})
	}
}

func TestDeleteSigner(t *testing.T) {
	now := time.Now()
	testTable := []struct {
		scenario      string
		expectedUser  *models.User
		expectedError error
	}{
		{
			"delete user",
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
			"can not delete user",
			&models.User{},
			&apperrors.CanNotDeleteUserErr,
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
			userRepoMock.EXPECT().DeleteUser(ctx, tc.expectedUser).Return(tc.expectedError)
			err := uInteractor.DeleteSigner(ctx, tc.expectedUser)
			if err != nil {

				if tc.expectedError != nil && apperrors.Is(err, tc.expectedError.(*apperrors.AppError)) {
					return
				}

				t.Fatal(err)
			}

		})
	}

}

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

func TestFindSigners(t *testing.T) {
	testTable := []struct {
		scenario           string
		expectedPagination *models.Pagination
		expectedUsers      []*models.User
		expectedError      error
	}{
		{
			"find user with limit=5",
			&models.Pagination{Limit: 5},
			make([]*models.User, 5),
			nil,
		},
		{
			"users not found",
			&models.Pagination{Limit: 5},
			nil,
			&apperrors.PaginationErr,
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

	for _, testCase := range testTable {
		t.Run(testCase.scenario, func(t *testing.T) {
			ctx := context.Background()
			userRepoMock.EXPECT().FindUsers(ctx, testCase.expectedPagination).Return(testCase.expectedPagination, testCase.expectedUsers, testCase.expectedError)
			pagination, users, err := uInteractor.FindSigners(ctx, testCase.expectedPagination)
			if err != nil {

				if testCase.expectedError != nil && apperrors.Is(err, testCase.expectedError.(*apperrors.AppError)) {
					return
				}

				t.Fatal(err)
			}
			assert.Equal(t, pagination.Limit, testCase.expectedPagination.Limit)
			assert.Equal(t, users, testCase.expectedUsers)

		})
	}
}

func TestNewUserInteractor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	testTable := []struct {
		scenario                string
		inputUserRepository     repository.UserRepository
		inputHashSalt           string
		inputSigningKey         []byte
		InputExpireDuration     int
		expectedUserInterfactor *userInteractor
	}{
		{
			"userInterfactor successfully created ",
			userRepoMock,
			"hash_aslt",
			[]byte("signing_key"),
			1,
			&userInteractor{
				userRepo:       userRepoMock,
				hashSalt:       "hash_aslt",
				signingKey:     []byte("signing_key"),
				expireDuration: 1,
			},
		},
		{
			"userRepository is absent",
			nil,
			"hash_aslt",
			[]byte("signing_key"),
			1,
			&userInteractor{
				userRepo:       nil,
				hashSalt:       "hash_aslt",
				signingKey:     []byte("signing_key"),
				expireDuration: 1,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.scenario, func(t *testing.T) {

			ui := NewUserInteractor(testCase.inputUserRepository, testCase.inputHashSalt, testCase.inputSigningKey, testCase.InputExpireDuration)
			assert.Equal(t, ui, testCase.expectedUserInterfactor)

		})

	}

}

func TestHashing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := &userInteractor{
		userRepo:       userRepoMock,
		hashSalt:       "hash_aslt",
		signingKey:     []byte("signing_key"),
		expireDuration: 0,
	}

	testingHashingFunc := func(password string) string {
		pwd := sha1.New()
		pwd.Write([]byte(password))
		pwd.Write([]byte(uInteractor.hashSalt))
		return fmt.Sprintf("%x", pwd.Sum(nil))
	}

	testTable := []struct {
		scenario               string
		inputPassword          string
		expectedHashedPassword string
		expectedError          error
	}{
		{
			"hashed password successfully created",
			"password",
			testingHashingFunc("password"),
			nil,
		},
		{
			"empty password field",
			"",
			testingHashingFunc(""),
			errors.New("empty pasword field"),
		},
		{
			"empty hash salt field",
			"password",
			testingHashingFunc("password"),
			errors.New("empty hashSalt field"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.scenario, func(t *testing.T) {

			if testCase.scenario == "empty hash salt field" {
				uInteractor.hashSalt = ""
			}

			hashedPassword, err := uInteractor.hashing(testCase.inputPassword)
			if err != nil {
				assert.Equal(t, err, testCase.expectedError)
				return
			}

			assert.Equal(t, hashedPassword, testCase.expectedHashedPassword)
		})

	}

}
