package api

import (
	"github.com/kurkop/gojobs/cmd/gojob-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func CronJobGroup(g *echo.Group) {
	g.POST("/", handlers.CreateCronJob)
	g.GET("/", handlers.GetAllCronJob)
	g.GET("/:name", handlers.GetCronJob)
	g.DELETE("/:name", handlers.DeleteCronJob)
}
