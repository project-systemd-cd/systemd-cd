package echo

import (
	errorss "errors"
	"net/http"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"

	"github.com/labstack/echo/v4"
)

func pipelinesNameJobsIdGet(c echo.Context) (err error) {
	name := c.Param("name")
	groupId := c.Param("id")

	logger.Logger().Trace("-----------------------------------------------------------")
	logger.Logger().Trace("START - GET /pipelines/:name/jobs/:id")
	logger.Logger().Tracef("< :name = %s", name)
	logger.Logger().Tracef("< RemoteAddr = %s", c.Request().RemoteAddr)
	logger.Logger().Trace("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Trace("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		notFound := errorss.As(err, &ErrNotFound)
		if err == nil || notFound {
			if notFound {
				err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			}
			logger.Logger().Tracef("> Status = %d", c.Response().Status)
			logger.Logger().Tracef("> ContentLength = %d", c.Response().Size)
			logger.Logger().Tracef("END    - GET /pipelines/:name/jobs/:id %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines/:name/jobs/:id")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Trace("-----------------------------------------------------------")
	}()

	_, err = CheckJWT(c)
	if err != nil {
		err = c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
		return err
	}

	p, err := service.FindPipeline(name)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			return err
		}
		return err
	}

	res, err := p.GetJob(groupId)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			return err
		}
		return err
	}

	if len(res) == 0 {
		err = c.NoContent(http.StatusNoContent)
		return err
	}
	err = c.JSONPretty(http.StatusOK, res, "	")
	return err
}
