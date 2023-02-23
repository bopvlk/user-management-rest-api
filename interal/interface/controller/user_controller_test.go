package controller

import (
	"context"
	"encoding/json"
	"git.foxminded.com.ua/3_REST_API/gen/mocks"
	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/domain/requests"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetOneUserHandler(t *testing.T) {

	now := time.Now()
	user := &models.User{
		ID:        1234,
		UserName:  "John",
		FirstName: "John",
		LastName:  "Hall",
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	userId := strconv.Itoa(int(user.ID))

	testTable := []struct {
		scenario     string
		expectedUser *models.User
		response     *requests.GetOneUserResponse
		httpCode     int
	}{
		{
			"get one user by id",
			user,
			mappers.MapUserToGetUserResponse(user),
			http.StatusOK,
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

			req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(userId)

			userRepoMock.EXPECT().FindOneUserByID(ctx, tc.expectedUser.ID).Return(tc.expectedUser, nil)
			err := uController.GetOneUserHandler(c)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tc.httpCode, rec.Code)

			marshalledResponse, err := json.Marshal(tc.response)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, string(marshalledResponse), strings.TrimSuffix(rec.Body.String(), "\n"))
		})
	}

}
