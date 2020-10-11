package router

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	jobGroup := e.Group("/api/v1/job")

	// set routes
	api.MainGroup(e)
	api.JobGroup(jobGroup)

	return e
}
