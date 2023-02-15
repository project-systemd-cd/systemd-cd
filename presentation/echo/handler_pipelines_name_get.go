package echo

import (
	"net/http"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/systemd"
	"time"

	errorss "errors"

	"github.com/labstack/echo/v4"
)

type QueryParamPipelineGet struct {
	Embed []string `query:"embed"`
	From  *string  `query:"from"`
	To    *string  `query:"to"`
	Asc   bool     `query:"asc"`
}

type ResPipelineGet struct {
	Name            string               `json:"name"`
	Status          string               `json:"status"`
	CommitRef       string               `json:"commit_ref"`
	SystemdServices *[]SystemdServiceGet `json:"systemd_services,omitempty"`
}

type ResPipelineGetJobsEmbed struct {
	ResPipelineGet
	Jobs [][]pipeline.Job `json:"jobs"`
}

type SystemdServiceGet struct {
	Name   string         `json:"name"`
	Status systemd.Status `json:"status"`
}

func pipelinesNameGet(c echo.Context) (err error) {
	name := c.Param("name")

	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - GET /pipelines/:name")
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
			logger.Logger().Infof("END    - GET /pipelines/:name %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines/:name")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	query := QueryParamPipelineGet{}
	if err = c.Bind(&query); err != nil {
		return err
	}
	query.Embed = c.QueryParams()["embed"]

	p, err := repository.FindPipeline(name)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		return err
	}
	status := p.GetStatus()
	commitRef := p.GetCommitRef()

	systemdServices := &[]SystemdServiceGet{}
	res := ResPipelineGet{
		Name:            name,
		Status:          string(status),
		CommitRef:       commitRef,
		SystemdServices: systemdServices,
	}

	services, err := p.GetStatusSystemdServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		*res.SystemdServices = append(*res.SystemdServices, SystemdServiceGet{Name: name, Status: s.Status})
	}

	embedJobs := false
	for _, embed := range query.Embed {
		if embed == "jobs" {
			embedJobs = true
		}
	}

	query2 := pipeline.QueryParamJob{}
	if embedJobs {
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

		jobs, err := p.GetJobs(query2)
		if err != nil {
			return err
		}

		return c.JSONPretty(http.StatusOK, ResPipelineGetJobsEmbed{res, jobs}, "	")
	} else {
		return c.JSONPretty(http.StatusOK, res, "	")
	}
}
