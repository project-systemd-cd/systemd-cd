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
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		notFound := errorss.As(err, &ErrNotFound)
		if err == nil || notFound {
			if notFound {
				err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			}
			logger.Logger().Debugf("> Status = %d", c.Response().Status)
			logger.Logger().Tracef("> ContentLength = %d", c.Response().Size)
			logger.Logger().Debugf("END    - GET /pipelines/:name/jobs %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines/:name/jobs")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	_, err = CheckJWT(c)
	if err != nil {
		err = c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
		return err
	}

	query := QueryParamPipelineJobsGet{}
	if err = c.Bind(&query); err != nil {
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

	query2 := pipeline.QueryParamJob{}
	if query.From != nil {
		var from time.Time
		from, err = parseTime(*query.From)
		if err != nil {
			err = c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
			return err
		}
		query2.From = &from
	}
	if query.To != nil {
		var to time.Time
		to, err = parseTime(*query.To)
		if err != nil {
			err = c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
			return err
		}
		query2.To = &to
	}
	query2.Asc = query.Asc

	res, err := p.GetJobs(query2)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		err = c.NoContent(http.StatusNoContent)
		return err
	}
	err = c.JSONPretty(http.StatusOK, res, "	")
	return err
}
