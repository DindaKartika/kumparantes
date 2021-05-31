package router

import (
	controller "kumparantes/controllers"
	"kumparantes/controllers/middlewares"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	middlewares.SetMiddlewares(e)

	// set main routes
	controller.MainRoutes(e)

	return e
}
