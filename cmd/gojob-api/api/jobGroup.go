package api

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func JobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateJob)
	g.GET("/:id", handlers.GetJob)
	g.PUT("/:id", handlers.UpdateJob)
	g.DELETE("/:id", handlers.DeleteJob)
}
