package api

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func CronJobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateCronJob)
	g.GET("/:namespace/:name", handlers.GetCronJob)
	g.PUT("/:name", handlers.UpdateCronJob)
	g.DELETE("/:namespace/:name", handlers.DeleteCronJob)
}
