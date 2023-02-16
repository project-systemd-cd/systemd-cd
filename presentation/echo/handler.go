package echo

import (
	"github.com/labstack/echo/v4"
)

func registerHandler(e *echo.Echo, jwtIssuer string, jwtSecret string) {
	e.POST("/users/sign_in", usersSigninPost)
	e.GET("/pipelines", pipelinesGet)
	e.GET("/pipelines/:name", pipelinesNameGet)
	e.GET("/pipelines/:name/jobs", pipelinesNameJobsGet)
	e.GET("/pipelines/:name/jobs/:id", pipelinesNameJobsIdGet)
}
