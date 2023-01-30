package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type userController struct {
	userInteractor interactor.UserInteractor
}

type UserController interface {
	SingUpHandler(ctx echo.Context) error
	GetOneUserHandler(ctx echo.Context) error
	GetAllUsersHandler(ctx echo.Context) error
}

type response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	IsError bool   `json:"is_error"`
}

func NewUserController(us interactor.UserInteractor) UserController {
	return &userController{us}
}

func (uc *userController) SingUpHandler(c echo.Context) error {
	var params models.User

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	duration, token, err := uc.userInteractor.SignUp(c.Request().Context(), &params)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.Expires = time.Now().Add(*duration)
	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, response{Token: token, Message: "You were logged in!"})
}

func (uc *userController) GetOneUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter the right ID", err.Error())
	}
	user := &models.User{ID: uint(id)}
	user, err = uc.userInteractor.FindSignerByID(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

///////////////////////////restructed zone /////////////////////////////////

func (uc *userController) GetAllUsersHandler(c echo.Context) error {
	name := getUserClaims(c).User.UserName

	var allUsers []*models.User

	allUsers, err := uc.userInteractor.FindAllSigners(c.Request().Context(), allUsers)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"welcome": fmt.Sprintf("Hello,%v U are in restricted zone", name),
		"users":   allUsers,
	})
}

func getUserClaims(c echo.Context) *interactor.AuthClaims {
	return c.Get("user").(*jwt.Token).Claims.(*interactor.AuthClaims)
}
