package router

import (
	"git.foxminded.com.ua/3_REST_API/interal/config"
	appMiddleware "git.foxminded.com.ua/3_REST_API/interal/infrastructure/middleware"
	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	v "git.foxminded.com.ua/3_REST_API/interal/validator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, config *config.Config, appController *controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &v.CustomValidator{Validator: validator.New()}

	apiGroup := e.Group("/api/v1")
	apiGroup.POST("/sing-up", appController.SignUpHandler)
	apiGroup.POST("/sing-in", appController.SignInHandler)

	restrictedGroup := apiGroup.Group("/restricted")
	restrictedGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(interactor.AuthClaims)
		},
		SigningKey:  []byte(config.SigningKey),
		TokenLookup: "cookie:Authorization",
	}))

	restrictedGroup.GET("/user/:id", appController.GetOneUserHandler)
	restrictedGroup.GET("/users", appController.GetUsersHandler, appMiddleware.ModeratorRoleMiddleware)
	restrictedGroup.DELETE("/user/:id", appController.DeleteUserHandler, appMiddleware.AdminRoleMiddleware)
	restrictedGroup.PUT("/user/:id", appController.UpdateUserHandler, appMiddleware.AdminRoleMiddleware)
	restrictedGroup.DELETE("/user/profile", appController.DeleteOwnerProfileHandler)
	restrictedGroup.PUT("/user/profile", appController.UpdateOwnerProfileHandler)
	restrictedGroup.PATCH("/user/:username", appController.RateUserHandler)

	return e
}
