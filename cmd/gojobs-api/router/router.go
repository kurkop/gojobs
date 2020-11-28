package router

import (
	"github.com/kurkop/gojobs/cmd/gojobs-api/api"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	jobGroup := e.Group("/api/v1/gojobs/jobs")
	cronJobGroup := e.Group("/api/v1/gojobs/cronjobs")

	// set routes
	api.MainGroup(e)
	api.JobGroup(jobGroup)
	api.CronJobGroup(cronJobGroup)

	return e
}
