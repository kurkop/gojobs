package middlewares

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetKeyAuthMiddlewares(g *echo.Group) {
	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		// Skipper:    DefaultSkipper,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("GOJOBS_API_KEY"), nil
		},
	}))
}
