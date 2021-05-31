package api

import (
	"kumparantes/controllers/handlers"

	"github.com/labstack/echo"
)

func MainRoutes(e *echo.Echo) {
	e.GET("/articles", handlers.GetArticles)
	e.POST("/articles", handlers.AddArticles)
}
