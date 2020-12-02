package router

import (
	"github.com/kurkop/gojobs/cmd/gojobs-api/api"
	"github.com/kurkop/gojobs/cmd/gojobs-api/api/middlewares"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	jobGroup := e.Group("/api/v1/gojobs/jobs")
	cronJobGroup := e.Group("/api/v1/gojobs/cronjobs")

	middlewares.SetKeyAuthMiddlewares(jobGroup)
	middlewares.SetKeyAuthMiddlewares(cronJobGroup)
	// set routes
	api.MainGroup(e)
	api.JobGroup(jobGroup)
	api.CronJobGroup(cronJobGroup)

	return e
}
