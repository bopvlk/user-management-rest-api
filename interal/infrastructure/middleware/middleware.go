package middleware

import (
	"errors"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/mappers"
	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	"github.com/labstack/echo/v4"
)

func AdminRoleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := controller.GetUserClaims(c)

		switch {
		case claims.User.Role == "user":
			appErr := apperrors.WrongRoleErr.AppendMessage(errors.New("you are have a role: 'user'"))
			c.Logger().Error(appErr.Error())
			return mappers.MapAppErrorToHTTPError(appErr)
		case claims.User.Role == "moderator":
			appErr := apperrors.WrongRoleErr.AppendMessage(errors.New("you are have a role: 'moderator'"))
			c.Logger().Error(appErr.Error())
			return mappers.MapAppErrorToHTTPError(appErr)
		}
		return next(c)
	}
}

func ModeratorRoleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := controller.GetUserClaims(c)

		switch {
		case claims.User.Role == "user":
			appErr := apperrors.WrongRoleErr.AppendMessage(errors.New("you are have a role: 'user'"))
			c.Logger().Error(appErr.Error())
			return mappers.MapAppErrorToHTTPError(appErr)
		}
		return next(c)
	}
}
