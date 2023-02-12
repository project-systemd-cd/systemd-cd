package echo

import "github.com/labstack/echo/v4"

func registerHandler(e *echo.Echo) {
	e.GET("/pipelines", pipelinesGet)
	e.GET("/pipelines/:name", pipelinesNameGet)
}
