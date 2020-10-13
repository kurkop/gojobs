package main

import (
	"github.com/kurkop/gojob/cmd/gojob-api/config"
	"github.com/kurkop/gojob/cmd/gojob-api/router"
	"github.com/labstack/echo/v4/middleware"
)

// middlewares

func main() {
	// Kubernetes client init
	config.KubeConnect()

	// Open router
	e := router.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Start("0.0.0.0:8000")
}
