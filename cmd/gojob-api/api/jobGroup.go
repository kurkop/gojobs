package api

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func JobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateJob)
	// list
	g.GET("/", handlers.GetAllJob)
	g.GET("/:name", handlers.GetJob)
	g.DELETE("/:name", handlers.DeleteJob)
	// TODO:
	// Replace
	// g.PUT("/:name", handlers.UpdateJob)
	// Patch Atomic changes
	// g.PATCH("/:name", handlers.UpdateJob)
}
