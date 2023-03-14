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
	apiGroup.POST("/sing-up", func(c echo.Context) error { return appController.SignUpHandler(c) })
	apiGroup.POST("/sing-in", func(c echo.Context) error { return appController.SignInHandler(c) })

	restrictedGroup := apiGroup.Group("/restricted")
	restrictedGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(interactor.AuthClaims)
		},
		SigningKey:  []byte(config.SigningKey),
		TokenLookup: "cookie:Authorization",
	}))

	restrictedGroup.GET("/user/:id", func(c echo.Context) error { return appController.GetOneUserHandler(c) })
	restrictedGroup.GET("/users", func(c echo.Context) error { return appController.GetUsersHandler(c) }, appMiddleware.ModeratorRoleMiddleware)
	restrictedGroup.DELETE("/user/:id", func(c echo.Context) error { return appController.DeleteUserHandler(c) }, appMiddleware.AdminRoleMiddleware)
	restrictedGroup.PUT("/user/:id", func(c echo.Context) error { return appController.UpdateUserHandler(c) }, appMiddleware.AdminRoleMiddleware)
	restrictedGroup.DELETE("/user/profile", func(c echo.Context) error { return appController.DeleteOwnerProfileHandler(c) })
	restrictedGroup.PUT("/user/profile", func(c echo.Context) error { return appController.UpdateOwnerProfileHandler(c) })
	restrictedGroup.POST("/user/:username", func(c echo.Context) error { return appController.RateHandler(c) })

	return e
}
