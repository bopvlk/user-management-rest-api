package router

import (
	"git.foxminded.com.ua/3_REST_API/interal/interface/controller"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	v1 := e.Group("/api/v1")

	v1.POST("/sing-up", func(context echo.Context) error { return c.User.SingUpHandler(context) })
	v1.GET("/user/:id", func(context echo.Context) error { return c.User.GetOneUserHandler(context) })

	r1 := v1.Group("/restricted")

	r1.Use(echojwt.JWT([]byte("secret")))

	r1.GET("/users", func(context echo.Context) error { return c.User.GetAllUsersHandler(context) })

	return e
}
