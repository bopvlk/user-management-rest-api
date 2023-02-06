package controller

import (
	"fmt"
	"net/http"
	"strconv"

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
	GetAllUsersHandler(ctx echo.Context) error
	SignInHandler(c echo.Context) error
	DeleteUserHandler(c echo.Context) error
}

type signUpResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	IsError bool   `json:"is_error"`
}

func NewUserController(us interactor.UserInteractor) UserController {
	return &userController{us}
}

func (uc *userController) SignUpHandler(c echo.Context) error {
	var params models.User

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	duration, token, err := uc.userInteractor.SignUp(c.Request().Context(), &params)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "change username and try again")
	}

	saveAuthcookie(c, token, int(duration.Seconds()))

	return c.JSON(http.StatusCreated, signUpResponse{Token: token, Message: "You were logged in!"})
}

func (uc *userController) SignInHandler(c echo.Context) error {
	var params models.User
	if err := c.Bind(&params); err != nil {
		c.Logger().Error("problem param", err)
		return echo.NewHTTPError(http.StatusBadRequest, "We got some problem with request param. ")

	}
	if params.Password == "" || params.UserName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "We got some problem with request param. ")
	}

	duration, token, err := uc.userInteractor.SignIn(c.Request().Context(), &params)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Please check out username and pasword, and try again")
	}

	saveAuthcookie(c, token, int(duration.Seconds()))

	return c.JSON(http.StatusOK, signUpResponse{Token: token, Message: "You were logged in!"})
}

func (uc *userController) GetOneUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter the right ID", err.Error())
	}
	user := &models.User{ID: uint(id)}
	user, err = uc.userInteractor.FindSignerByID(c.Request().Context(), user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "some problem with request")
	}

	return c.JSON(http.StatusOK, user)
}

///////////////////////////restructed zone /////////////////////////////////

func (uc *userController) GetAllUsersHandler(c echo.Context) error {
	name := getUserClaims(c).User.UserName

	var allUsers []*models.User

	allUsers, err := uc.userInteractor.FindAllSigners(c.Request().Context(), allUsers)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "some problem with request handling")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"welcome": fmt.Sprintf("Hello,%v U are in restricted zone", name),
		"users":   allUsers,
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
