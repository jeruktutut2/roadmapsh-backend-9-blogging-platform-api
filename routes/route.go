package routes

import (
	"blogging-platform-api/controllers"

	"github.com/labstack/echo/v4"
)

func BlogRoute(e *echo.Echo, controller controllers.BlogController) {
	e.POST("/posts", controller.Create)
	e.PUT("/posts/:id", controller.Update)
	e.DELETE("/posts/:id", controller.Delete)
	e.GET("/posts/:id", controller.FindById)
	e.GET("/posts", controller.FindAll)
}
