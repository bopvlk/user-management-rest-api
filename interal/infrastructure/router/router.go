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

func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	v1 := e.Group("/api/v1")

	v1.POST("/sing-up", func(context echo.Context) error { return c.User.SingUpHandler(context) })
	v1.GET("/user/:id", func(context echo.Context) error { return c.User.GetOneUserHandler(context) })

	r1 := v1.Group("/restricted")

	r1.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(context echo.Context) jwt.Claims {
			return new(interactor.AuthClaims)
		},
		SigningKey:  []byte(viper.GetString("SIGNING_KEY")),
		TokenLookup: "cookie:Authorization",
	}))

	r1.GET("/users", func(context echo.Context) error { return c.User.GetAllUsersHandler(context) })

	return e
}
