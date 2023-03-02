package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/domain/requests"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type userController struct {
	userInteractor interactor.UserInteractor
}

type UserController interface {
	SignUpHandler(ctx echo.Context) error
	GetOneUserHandler(ctx echo.Context) error
	GetUsersHandler(ctx echo.Context) error
	SignInHandler(c echo.Context) error
	DeleteUserHandler(c echo.Context) error
	UpdateUserHandler(c echo.Context) error
}

func NewUserController(us interactor.UserInteractor) UserController {
	return &userController{us}
}

func (uc *userController) SignUpHandler(c echo.Context) error {
	var signUpRequest requests.SignUpRequest
	if err := c.Bind(&signUpRequest); err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	if err := c.Validate(signUpRequest); err != nil {
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	duration, token, err := uc.userInteractor.SignUp(c.Request().Context(), mappers.MapSignUpRequestToUser(&signUpRequest))
	if err != nil {
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	saveAuthcookie(c, token, duration)

	return c.JSON(http.StatusCreated, requests.SignUpInResponse{Token: token, Message: "You are logged in!"})
}

func (uc *userController) SignInHandler(c echo.Context) error {
	var signInRequest requests.SignInRequest
	if err := c.Bind(&signInRequest); err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	duration, token, err := uc.userInteractor.SignIn(c.Request().Context(), signInRequest.UserName, signInRequest.Password)
	if err != nil {
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	saveAuthcookie(c, token, duration)

	return c.JSON(http.StatusOK, requests.SignUpInResponse{Token: token, Message: "You are logged in!"})
}

func (uc *userController) GetOneUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	user, err := uc.userInteractor.FindOneSigner(c.Request().Context(), uint(id))
	if err != nil {
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	return c.JSON(http.StatusOK, mappers.MapUserToGetUserResponse(user))
}

func (uc *userController) GetUsersHandler(c echo.Context) error {
	if GetUserClaims(c).User.Role == "user" {
		appErr := &apperrors.WrongRoleErr
		c.Logger().Error(appErr.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	name := GetUserClaims(c).User.UserName
	pagination := mappers.MapContextToPagination(c)
	pagination, users, err := uc.userInteractor.FindSigners(c.Request().Context(), pagination)
	if err != nil {
		c.Logger().Error(err)
		return mappers.MapAppErrorToHTTPError(err)
	}

	urlPath := c.Request().URL.Path
	pagination.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 1, pagination.Sort)
	pagination.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.TotalPages, pagination.Sort)

	if pagination.Page > 1 {
		pagination.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.Page-1, pagination.Sort)
	}

	if pagination.Page < pagination.TotalPages {
		pagination.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.Page+1, pagination.Sort)
	}

	return c.JSON(http.StatusOK, mappers.MapPaginationAndUsersToGetUsersResponse(users, pagination, name))
}

func (uc *userController) DeleteUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	if err := uc.userInteractor.DeleteSignerByID(c.Request().Context(), id); err != nil {
		c.Logger().Warn(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("User with id:%d is deleted", id))
}

func (uc *userController) UpdateUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	claims := GetUserClaims(c)
	switch {
	case claims.User.Role == "user" && claims.User.ID != uint(id):
		appErr := apperrors.WrongRoleErr.AppendMessage(errors.New("you are have a role: 'user'"))
		c.Logger().Error(appErr.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	case claims.User.Role == "moderator" && claims.User.ID != uint(id):
		appErr := apperrors.WrongRoleErr.AppendMessage(errors.New("you are have a role: 'moderator'"))
		c.Logger().Error(appErr.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	var updateRequest requests.UpdateRequest
	if err := c.Bind(&updateRequest); err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(err.Error())
		return mappers.MapAppErrorToHTTPError(appErr)
	}

	user, err := uc.userInteractor.UpdateSignersByID(c.Request().Context(), id, int(claims.User.ID), mappers.MapUpdateRequestToUser(&updateRequest))
	if err != nil {
		c.Logger().Warn(err.Error())
		return mappers.MapAppErrorToHTTPError(err)
	}

	return c.JSON(http.StatusOK, mappers.MapUserToUpdateResponse(user))
}

func GetUserClaims(c echo.Context) *interactor.AuthClaims {
	return c.Get("user").(*jwt.Token).Claims.(*interactor.AuthClaims)
}

func saveAuthcookie(c echo.Context, token string, duration int) {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.MaxAge = duration
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}
