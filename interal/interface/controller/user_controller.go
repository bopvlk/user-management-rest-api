package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
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
}

func NewUserController(us interactor.UserInteractor) UserController {
	return &userController{us}
}

func (uc *userController) SignUpHandler(c echo.Context) error {
	var signUpRequest requests.SignUpRequest
	if err := c.Bind(&signUpRequest); err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(appErr.Code)
		errResponse := mappers.MapAppErrorToErrorResponse(appErr)
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	if err := c.Validate(signUpRequest); err != nil {
		appErr := apperrors.ValidatorErr.AppendMessage(err)
		c.Logger().Error(appErr.Code)
		errResponse := mappers.MapAppErrorToErrorResponse(appErr)
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	duration, token, err := uc.userInteractor.SignUp(c.Request().Context(), mappers.MapSignUpRequestToUser(&signUpRequest))
	if err != nil {
		c.Logger().Error(err.(*apperrors.AppError).Code)
		errResponse := mappers.MapAppErrorToErrorResponse(err.(*apperrors.AppError))
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	saveAuthcookie(c, token, duration)

	return c.JSON(http.StatusCreated, requests.SignUpInResponse{Token: token, Message: "You are logged in!"})
}

func (uc *userController) SignInHandler(c echo.Context) error {
	var signInRequest requests.SignInRequest
	if err := c.Bind(&signInRequest); err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(appErr.Code)
		errResponse := mappers.MapAppErrorToErrorResponse(appErr)
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	if err := c.Validate(signInRequest); err != nil {
		appErr := apperrors.ValidatorErr.AppendMessage(err)
		c.Logger().Error(appErr.Code)
		errResponse := mappers.MapAppErrorToErrorResponse(appErr)
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	if signInRequest.Password == "" || signInRequest.UserName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "We got some problem with request param. ")
	}

	duration, token, err := uc.userInteractor.SignIn(c.Request().Context(), signInRequest.UserName, signInRequest.Password)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Please check out username and pasword, and try again")
	}

	saveAuthcookie(c, token, duration)

	return c.JSON(http.StatusOK, requests.SignUpInResponse{Token: token, Message: "You were logged in!"})
}

func (uc *userController) GetOneUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appErr := apperrors.CanNotBindErr.AppendMessage(err)
		c.Logger().Error(appErr.Code)
		errResponse := mappers.MapAppErrorToErrorResponse(appErr)
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	user, err := uc.userInteractor.FindOneSigner(c.Request().Context(), uint(id))
	if err != nil {
		c.Logger().Error(err.(*apperrors.AppError).Code)
		errResponse := mappers.MapAppErrorToErrorResponse(err.(*apperrors.AppError))
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	return c.JSON(http.StatusOK, requests.GetOneUserResponse{
		Message:      fmt.Sprintf("There is user with ID %v", id),
		UserResponse: *mappers.MapUserToUserResponse(user),
		IsError:      false,
	})
}

func (uc *userController) GetUsersHandler(c echo.Context) error {
	name := getUserClaims(c).User.UserName
	pagination, err := generatePaginationRequest(c)
	if err != nil {
		c.Logger().Error(err.(*apperrors.AppError).Code)
		errResponse := mappers.MapAppErrorToErrorResponse(err.(*apperrors.AppError))
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	pagination, users, err := uc.userInteractor.FindSigners(c.Request().Context(), pagination)
	if err != nil {
		c.Logger().Error(err.(*apperrors.AppError).Code)
		errResponse := mappers.MapAppErrorToErrorResponse(err.(*apperrors.AppError))
		return echo.NewHTTPError(errResponse.HTTPCode, errResponse)
	}

	pagination.Rows = mappers.MapUsersToUsersResponse(users)

	urlPath := c.Request().URL.Path

	pagination.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 1, pagination.Sort)
	pagination.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.TotalPages, pagination.Sort)

	if pagination.Page > 1 {
		pagination.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.Page-1, pagination.Sort)
	}

	if pagination.Page < pagination.TotalPages {
		pagination.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, pagination.Page+1, pagination.Sort)
	}

	return c.JSON(http.StatusOK, requests.GetUsersResponse{
		Message:       fmt.Sprintf("Hello,%v U are in restricted zone", name),
		UsersResponse: pagination,
	})
}

func (uc *userController) DeleteUserHandler(c echo.Context) error {
	user := getUserClaims(c).User

	if err := uc.userInteractor.DeleteSigner(c.Request().Context(), user); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "some problem with request handling")
	}
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "some problem with request handling")
	}
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, fmt.Sprintf("User %s with id:%d is deleted", user.UserName, user.ID))
}

func getUserClaims(c echo.Context) *interactor.AuthClaims {
	return c.Get("user").(*jwt.Token).Claims.(*interactor.AuthClaims)
}

func saveAuthcookie(c echo.Context, token string, duration int) {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.MaxAge = duration
	c.SetCookie(cookie)
}

func generatePaginationRequest(c echo.Context) (*models.Pagination, error) {
	var err error
	limit, errLimit := strconv.Atoi(c.QueryParam("limit"))
	if errLimit != nil {
		err = errLimit
		limit = 5
	} else if limit < 5 {
		limit = 5
	}

	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		err = fmt.Errorf("%v :  %v", err, errPage)
		if err != nil {
			page = 1
		} else if page < 1 {
			page = 1
		}
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		err = fmt.Errorf("%v : query param 'sort' is not correct", err)
		sort = "id desc"
	}

	return &models.Pagination{Limit: limit, Page: page, Sort: sort}, apperrors.PaginationErr.AppendMessage(err)
}
