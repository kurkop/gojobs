package router

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	// TODO: Resource Namespace
	jobGroup := e.Group("/api/v1/:namespace/jobs")
	cronJobGroup := e.Group("/api/v1/:namespace/cronjobs")

	// set routes
	api.MainGroup(e)
	api.JobGroup(jobGroup)
	api.CronJobGroup(cronJobGroup)

	return e
}
