package api

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func JobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateJob)
	g.GET("/:namespace/:name", handlers.GetJob)
	g.PUT("/:name", handlers.UpdateJob)
	g.DELETE("/:namespace/:name", handlers.DeleteJob)
}
