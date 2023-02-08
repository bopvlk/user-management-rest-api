package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
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
	var params mappers.SignUpRequest
	if err := c.Bind(&params); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, mappers.SignUpInResponse{Message: "We got some problem with request param", IsError: true})
	}

	duration, token, err := uc.userInteractor.SignUp(c.Request().Context(), mappers.MapSignUpRequestToUserModel(&params))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusConflict, mappers.SignUpInResponse{Message: "change username and try again", IsError: true})
	}

	saveAuthcookie(c, token, int(duration.Seconds()))

	return c.JSON(http.StatusCreated, mappers.SignUpInResponse{Token: token, Message: "You were logged in!"})
}

func (uc *userController) SignInHandler(c echo.Context) error {
	var params mappers.SignInRequest
	if err := c.Bind(&params); err != nil {
		c.Logger().Error("problem param", err)
		return echo.NewHTTPError(http.StatusBadRequest, "We got some problem with request param.")
	}

	if params.Password == "" || params.UserName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "We got some problem with request param. ")
	}

	duration, token, err := uc.userInteractor.SignIn(c.Request().Context(), mappers.MapSignInRequestToUserModel(&params))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Please check out username and pasword, and try again")
	}

	saveAuthcookie(c, token, int(duration.Seconds()))

	return c.JSON(http.StatusOK, mappers.SignUpInResponse{Token: token, Message: "You were logged in!"})
}

func (uc *userController) GetOneUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, mappers.GetOneUserResponse{
			Message: "Please enter the correct ID",
			IsError: true,
		})
	}
	user := &models.User{ID: uint(id)}
	user, err = uc.userInteractor.FindOneSigner(c.Request().Context(), user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, mappers.GetOneUserResponse{
			Message: "some problem with request",
			IsError: true,
		})
	}

	return c.JSON(http.StatusOK, mappers.GetOneUserResponse{
		Message:      fmt.Sprintf("There is user with ID %v", id),
		UserResponse: *mappers.MapUserModelToUserResponse(user),
		IsError:      false,
	})
}

func (uc *userController) GetUsersHandler(c echo.Context) error {

	page, err := strconv.Atoi(c.Param("page"))

	name := getUserClaims(c).User.UserName

	var users []*models.User
	users, err = uc.userInteractor.FindSigners(c.Request().Context(), page, users)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "some problem with request handling")
	}

	return c.JSON(http.StatusOK, mappers.GetUsersResponse{
		Message:       fmt.Sprintf("Hello,%v U are in restricted zone", name),
		UsersResponse: mappers.MapUsersModelToUsersResponse(users),
		IsError:       false,
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
