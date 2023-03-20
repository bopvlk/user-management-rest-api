package controller

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.foxminded.com.ua/3_REST_API/gen/mocks"
	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/domain/requests"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	v "git.foxminded.com.ua/3_REST_API/interal/validator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandler(t *testing.T) {

	inputUser := getTestUser()
	inputUser.ID = 0
	inputUser.Rating.Rating = 1
	inputUser.Password = hashingUserFunc(inputUser.Password)

	testTable := []struct {
		scenario         string
		inputuser        *models.User
		expectedUser     *models.User
		signUpRequest    string
		expectedResponse requests.SignUpInResponse
		expectedhttpCode int
		expectedError    error
	}{
		{
			"user successfully redistered",
			inputUser,
			getTestUser(),
			`{"user_name": "JohnHall", "role": "admin", "first_name": "John", "last_name": "Hall", "password": "very12difficult()Password"}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusCreated,
			nil,
		},
		{
			"wrong request body",
			inputUser,
			getTestUser(),
			`{"use, "first_name": "John",  "Hall", "password": "}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"can not tgrough out a validation",
			inputUser,
			getTestUser(),
			`{"user_name": "JohnHall", "first_name": "John", "last_name": "Hall", "password": "1234"}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"tries to create a user with an existing username",
			inputUser,
			nil,
			`{"user_name": "JohnHall", "first_name": "John", "last_name": "Hall", "password": "very12difficult()Password"}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusBadRequest,
			&apperrors.CanNotCreateUserErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()
			e.Validator = &v.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/sing-up", strings.NewReader(tc.signUpRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			userRepoMock.EXPECT().CreateUser(ctx, tc.inputuser).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			err := uController.SignUpHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.expectedhttpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.expectedhttpCode, rec.Code)
			marshalledResponse, err := json.Marshal(tc.expectedResponse)

			if assert.NoError(t, err) {
				assert.Equal(t, string(marshalledResponse), strings.TrimSuffix(rec.Body.String(), "\n"))
			}

		})
	}
}

func TestSignInHandler(t *testing.T) {

	testTable := []struct {
		scenario         string
		inputuser        *models.User
		expectedUser     *models.User
		signUpRequest    string
		expectedResponse requests.SignUpInResponse
		expectedhttpCode int
		expectedError    error
	}{
		{
			"user successfully redistered",
			&models.User{
				ID:        0,
				UserName:  "JohnHall",
				Role:      "admin",
				FirstName: "John",
				LastName:  "Hall",
				Password:  hashingUserFunc("very12difficult()Password"),
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			getTestUser(),
			`{"user_name": "JohnHall", "password": "very12difficult()Password"}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusOK,
			nil,
		},
		{
			"wrong request body",
			&models.User{
				ID:        0,
				UserName:  "JohnHall",
				FirstName: "John",
				LastName:  "Hall",
				Password:  hashingUserFunc("very12difficult()Password"),
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			getTestUser(),
			`{"use, "first_name": "John",  "Hall", "password": "}`,
			requests.SignUpInResponse{},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"tries sign in with an wrong UserName or Password",
			&models.User{
				ID:        0,
				UserName:  "JohnHall",
				FirstName: "John",
				LastName:  "Hall",
				Password:  hashingUserFunc("very12difficult()Password"),
				CreatedAt: nil,
				UpdatedAt: nil,
			},
			nil,
			`{"user_name": "JohnHall", "password": "very12difficult()Password"}`,
			requests.SignUpInResponse{Message: "You are logged in!"},
			http.StatusBadRequest,
			&apperrors.UserNotFoundErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()
			e.Validator = &v.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/sing-in", strings.NewReader(tc.signUpRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			userRepoMock.EXPECT().FindOneUserByUserNameAndPassword(ctx, tc.inputuser.UserName, tc.inputuser.Password).
				Return(tc.expectedUser, tc.expectedError).AnyTimes()

			err := uController.SignInHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.expectedhttpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.expectedhttpCode, rec.Code)
			marshalledResponse, err := json.Marshal(tc.expectedResponse)

			if assert.NoError(t, err) {
				assert.Equal(t, string(marshalledResponse), strings.TrimSuffix(rec.Body.String(), "\n"))
			}

		})
	}
}

func TestGetOneUserHandler(t *testing.T) {

	user := getTestUser()
	userId := strconv.Itoa(int(user.ID))

	testTable := []struct {
		scenario      string
		expectedUser  *models.User
		inputUserID   string
		response      *requests.GetOneUserResponse
		httpCode      int
		expectedError error
	}{
		{
			"get one user by id",
			user,
			userId,
			mappers.MapUserToGetUserResponse(user),
			http.StatusOK,
			nil,
		},
		{
			"wrong path params",
			user,
			"userID",
			mappers.MapUserToGetUserResponse(user),
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"user not found by id",
			user,
			userId,
			mappers.MapUserToGetUserResponse(user),
			http.StatusBadRequest,
			&apperrors.UserNotFoundErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")

			c.SetParamValues(tc.inputUserID)
			userRepoMock.EXPECT().FindOneUserByID(ctx, tc.expectedUser.ID).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			err := uController.GetOneUserHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)
			marshalledResponse, err := json.Marshal(tc.response)

			if assert.NoError(t, err) {
				assert.Equal(t, string(marshalledResponse), strings.TrimSuffix(rec.Body.String(), "\n"))
			}

		})
	}
}

func TestGetUsersHandler(t *testing.T) {
	now := time.Now()

	expectedPagination := &models.Pagination{
		Limit: 5,
		Page:  1,
		Sort:  "id desc",
	}

	paginationLimit := strconv.Itoa(int(expectedPagination.Limit))
	paginationPage := strconv.Itoa(int(expectedPagination.Page))

	usersGenerator := func(pagination *models.Pagination) []*models.User {
		users := make([]*models.User, pagination.Limit)

		for i := 0; i < len(users); i++ {
			users[i] = &models.User{}
			users[i].ID = uint(((pagination.Page - 1) * pagination.Limit) + i + 1)
			users[i].UserName = fmt.Sprintf("UniqueUser%d", int(users[i].ID))
			users[i].Role = "admin"
			users[i].FirstName = fmt.Sprintf("FirstName%d", int(users[i].ID))
			users[i].LastName = fmt.Sprintf("LastName%d", int(users[i].ID))
			users[i].CreatedAt = &now
			users[i].UpdatedAt = &now
		}

		return users
	}

	fineUsers := usersGenerator(expectedPagination)

	tokenGenerator := func(user *models.User) *jwt.Token {
		claims := &interactor.AuthClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * (time.Duration(1)))),
			},
		}

		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	}

	testTable := []struct {
		scenario           string
		expectedPagination *models.Pagination
		expectedUsers      []*models.User
		httpCode           int
		expectedError      error
	}{
		{
			"get users",
			expectedPagination,
			fineUsers,
			http.StatusOK,
			nil,
		},
		{
			"wrong query param",
			&models.Pagination{
				Limit: 5,
				Page:  1,
				Sort:  "id desc",
			},
			fineUsers,
			http.StatusOK,
			nil,
		},
		{
			"bad user role",
			&models.Pagination{
				Limit: 5,
				Page:  1,
				Sort:  "id desc",
			},
			fineUsers,
			http.StatusForbidden,
			&apperrors.WrongRoleErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()
			q := make(url.Values)

			if tc.scenario != "wrong query param" {
				q.Set("limit", paginationLimit)
				q.Set("page", paginationPage)
				q.Set("sort", tc.expectedPagination.Sort)
			}

			userRepoMock.EXPECT().FindUsers(ctx, tc.expectedPagination).Return(tc.expectedPagination, tc.expectedUsers, tc.expectedError).AnyTimes()

			req := httptest.NewRequest(http.MethodGet, "/users?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.scenario == "bad user role" {
				tc.expectedUsers[0].Role = "user"
			}

			c.Set("user", tokenGenerator(tc.expectedUsers[0]))

			err := uController.GetUsersHandler(c)
			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

			var pag requests.GetUsersResponse
			json.NewDecoder(rec.Body).Decode(&pag)

			if assert.NoError(t, err) {
				for i, v := range pag.UsersResponse.Rows.([]interface{}) {
					assert.Equal(t, v.(map[string]interface{})["user_name"].(string), tc.expectedUsers[i].UserName)
					assert.Equal(t, v.(map[string]interface{})["first_name"].(string), tc.expectedUsers[i].FirstName)
					assert.Equal(t, v.(map[string]interface{})["last_name"].(string), tc.expectedUsers[i].LastName)
				}
			}

		})
	}
}

func TestDeleteUserHandler(t *testing.T) {

	expectedUserWithRoleUser := getTestUser()
	expectedUserWithRoleUser.Role = "user"

	testTable := []struct {
		scenario      string
		expectedID    string
		expectedUser  *models.User
		httpCode      int
		expectedError error
	}{
		{
			"successfully deleted user",
			"1234",
			getTestUser(),
			http.StatusOK,
			nil,
		},
		{
			"wrong path params",
			"userID",
			getTestUser(),
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"id is exsist",
			"12598",
			getTestUser(),
			http.StatusInternalServerError,
			&apperrors.CanNotDeleteUserErr,
		},
		{
			"role have lack of rights",
			"1234",
			expectedUserWithRoleUser,
			http.StatusForbidden,
			&apperrors.WrongRoleErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/users/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("id")
			c.SetParamValues(tc.expectedID)
			c.Set("user", tokenGenerator())

			id, _ := strconv.Atoi(tc.expectedID)
			userRepoMock.EXPECT().DeleteUserByID(ctx, id).Return(tc.expectedError).AnyTimes()

			err := uController.DeleteUserHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

		})
	}
}

func TestDeleteOwnerProfileHandler(t *testing.T) {

	testTable := []struct {
		scenario      string
		expectedID    int
		httpCode      int
		expectedError error
	}{
		{
			"successfully deleted user",
			124,
			http.StatusOK,
			nil,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
	uController := NewUserController(uInteractor)

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/user/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			userRepoMock.EXPECT().DeleteOwnUser(ctx, tc.expectedID).Return(tc.expectedError)
			c.Set("user", tokenGenerator())
			err := uController.DeleteOwnerProfileHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

		})
	}
}

func TestUpdateUserHandler(t *testing.T) {

	testTable := []struct {
		scenario              string
		expectedID            string
		expectedUpdateRequest string
		expectedUser          *models.User
		httpCode              int
		expectedError         error
	}{
		{
			"successfully updated user",
			"1234",
			`{"user_name": "JohnHall", "role": "user", "first_name": "John", "last_name": "Hall"}`,
			&models.User{
				ID:        0,
				UserName:  "JohnHall",
				FirstName: "John",
				Role:      "user",
				LastName:  "Hall",
			},
			http.StatusOK,
			nil,
		},
		{
			"wrong path params",
			"userID",
			`{}`,
			&models.User{},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"wrong bind params",
			"1234",
			`{"user_namsdsdfe": "JohnHal, "passdfsword": "very12difficult()Password"}`,
			&models.User{UserName: "JohnHall"},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"user has admin status",
			"1234",
			`{"user_name": "JohnHall", "role": "admin", "first_name": "John", "last_name": "Hall"}`,
			&models.User{
				ID:        0,
				UserName:  "JohnHall",
				FirstName: "John",
				Role:      "admin",
				LastName:  "Hall",
			},
			http.StatusInternalServerError,
			&apperrors.CanNotUpdateErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
	uController := NewUserController(uInteractor)

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPut, "/users/:id", strings.NewReader(tc.expectedUpdateRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("id")
			c.SetParamValues(tc.expectedID)

			id, _ := strconv.Atoi(tc.expectedID)

			userRepoMock.EXPECT().UpdateUserByID(ctx, id, tc.expectedUser).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			err := uController.UpdateUserHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

		})
	}
}

func TestUpdateOwnerProfileHandler(t *testing.T) {

	inputUser := getTestUser()
	inputUser.ID = 0
	inputUser.Rating.Rating = 0
	testTable := []struct {
		scenario string

		expectedUpdateRequest string
		expectedUser          *models.User
		inputUser             *models.User
		httpCode              int
		expectedError         error
	}{
		{
			"successfully updated profile",
			`{"user_name": "JohnHall", "role": "admin", "first_name": "John", "last_name": "Hall", "password": "very12difficult()Password"}`,
			getTestUser(),
			inputUser,
			http.StatusOK,
			nil,
		},

		{
			"wrong bind params",
			`{"user_name": ", "last_name": "Hall", "passwordy12difficult()Password"}`,
			getTestUser(),
			inputUser,
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := mocks.NewMockUserRepository(ctrl)
	uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
	uController := NewUserController(uInteractor)

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPut, "/user/profile", strings.NewReader(tc.expectedUpdateRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", tokenGenerator())

			userRepoMock.EXPECT().UpdateOwnUser(ctx, int(tc.expectedUser.ID), tc.inputUser).Return(tc.expectedUser, tc.expectedError).AnyTimes()

			err := uController.UpdateOwnerProfileHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

		})
	}
}

func TestRateHandler(t *testing.T) {

	expectedUser := getTestUser()
	expectedUser.UserName = "username"

	type args struct {
		ctx                 context.Context
		myID                string
		username            string
		expectedRatedUpDown bool
	}

	testTable := []struct {
		scenario              string
		expectedUpdateRequest string
		expectedUser          *models.User
		args                  args
		httpCode              int
		expectedError         error
	}{
		{
			"successfully rate profile",
			`{"rate": true}`,
			expectedUser,
			args{
				context.Background(),
				"124+",
				expectedUser.UserName,
				true,
			},
			http.StatusOK,
			nil,
		},

		{
			"wrong bind params",
			`{"rate": ", "last_name": "Hall", "passwordy12difficult()Password"}`,
			expectedUser,
			args{
				context.Background(),
				"124+",
				expectedUser.UserName,
				true,
			},
			http.StatusBadRequest,
			&apperrors.CanNotBindErr,
		},
		{
			"can not rate himself",
			`{"rate": true}`,
			getTestUser(),
			args{
				context.Background(),
				"124+",
				"JohnHall",
				true,
			},
			http.StatusForbidden,
			&apperrors.CanNotRateYorself,
		},
		// {
		// 	"missing particular user",
		// 	`{"rate": true}`,
		// 	expectedUser,
		// 	"fedir",
		// 	true,
		// 	http.StatusBadRequest,
		// 	&apperrors.UserNotFoundErr,
		// },
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range testTable {
		t.Run(tc.scenario, func(t *testing.T) {
			userRepoMock := mocks.NewMockUserRepository(ctrl)
			uInteractor := interactor.NewUserInteractor(userRepoMock, "hash_salt", []byte("signing_key"), 1)
			uController := NewUserController(uInteractor)

			e := echo.New()

			req := httptest.NewRequest(http.MethodPut, "/user/:username", strings.NewReader(tc.expectedUpdateRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("username")

			c.SetParamValues(tc.args.username)
			c.Set("user", tokenGenerator())

			userRepoMock.EXPECT().RateUserByUsername(ctx, tc.args.myID, tc.args.username, tc.args.expectedRatedUpDown).Return(tc.expectedUser, tc.expectedError).
				AnyTimes()

			err := uController.RateUserHandler(c)

			if err != nil {
				apperrors.Is(err, tc.expectedError.(*apperrors.AppError))
				assert.Equal(t, tc.httpCode, err.(*echo.HTTPError).Code)
				return
			}
			assert.Equal(t, tc.httpCode, rec.Code)

		})
	}
}

func hashingUserFunc(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte("hash_salt"))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}

func getTestUser() *models.User {

	return &models.User{
		ID:        124,
		UserName:  "JohnHall",
		Role:      "admin",
		Rating:    models.Rating{Rating: 10},
		FirstName: "John",
		LastName:  "Hall",
		Password:  "very12difficult()Password",
	}
}

func tokenGenerator() *jwt.Token {
	claims := &interactor.AuthClaims{
		User: getTestUser(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * (time.Duration(1)))),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
