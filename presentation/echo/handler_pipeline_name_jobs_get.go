package echo

import (
	errorss "errors"
	"net/http"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"

	"github.com/labstack/echo/v4"
)

type QueryParamPipelineJobsGet struct {
	From *string `query:"from"`
	To   *string `query:"to"`
	Asc  bool    `query:"asc"`
}

func pipelinesNameJobsGet(c echo.Context) (err error) {
	name := c.Param("name")

	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - GET /pipelines/:name/jobs")
	logger.Logger().Tracef("< :name = %s", name)
	logger.Logger().Tracef("< RemoteAddr = %s", c.Request().RemoteAddr)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		notFound := errorss.As(err, &ErrNotFound)
		if err == nil || notFound {
			if notFound {
				err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			}
			logger.Logger().Debugf("> Status = %d", c.Response().Status)
			logger.Logger().Tracef("> ContentLength = %d", c.Response().Size)
			logger.Logger().Infof("END    - GET /pipelines/:name/jobs %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines/:name/jobs")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	query := QueryParamPipelineJobsGet{}
	if err = c.Bind(&query); err != nil {
		return err
	}

	p, err := repository.FindPipeline(name)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		return err
	}

	query2 := pipeline.QueryParamJob{}
	if query.From != nil {
		var from time.Time
		from, err = time.Parse(time.RFC3339, *query.From)
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		query2.From = &from
	}
	if query.To != nil {
		var to time.Time
		to, err = time.Parse(time.RFC3339, *query.To)
		if err != nil {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		query2.To = &to
	}
	query2.Asc = query.Asc

	res, err := p.GetJobs(query2)
	if err != nil {
		return err
	}

	if res == nil {
		err = c.NoContent(http.StatusNoContent)
		return err
	}
	err = c.JSONPretty(http.StatusOK, res, "	")
	return err
}
