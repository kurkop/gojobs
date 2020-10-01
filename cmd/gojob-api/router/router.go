package router

import (
	"github.com/kurkop/gojob/cmd/gojob-api/api"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	// set routes
	api.MainGroup(e)

	return e
}
