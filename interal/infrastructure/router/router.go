package router

import (
	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/interactor"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func NewRouter(e *echo.Echo, appController *controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	apiGroup := e.Group("/api/v1")
	apiGroup.POST("/sing-up", func(c echo.Context) error { return appController.SignUpHandler(c) })
	apiGroup.GET("/user/:id", func(c echo.Context) error { return appController.GetOneUserHandler(c) })
	apiGroup.POST("/sing-in", func(c echo.Context) error { return appController.SignInHandler(c) })

	restrictedGroup := apiGroup.Group("/restricted")
	restrictedGroup.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(interactor.AuthClaims)
		},
		SigningKey:  []byte(viper.GetString("SIGNING_KEY")),
		TokenLookup: "cookie:Authorization",
	}))
	restrictedGroup.GET("/users", func(c echo.Context) error { return appController.GetUsersHandler(c) })
	restrictedGroup.DELETE("/users", func(c echo.Context) error { return appController.DeleteUserHandler(c) })
	return e
}
