package api

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func CronJobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateCronJob)
	g.GET("/:name", handlers.GetCronJob)
	g.DELETE("/:name", handlers.DeleteCronJob)
}
