package api

import (
	"github.com/kurkop/gojobs/cmd/gojobs-api/api/handlers"
	"github.com/labstack/echo/v4"
)

func MainGroup(e *echo.Echo) {
	e.GET("/", handlers.Hello)
}
