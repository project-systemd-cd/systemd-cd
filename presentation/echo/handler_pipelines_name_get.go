package echo

import (
	"net/http"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"

	errorss "errors"

	"github.com/labstack/echo/v4"
)

type ResPipelineGet struct {
	Name            string               `json:"name"`
	Status          string               `json:"status"`
	CommitRef       string               `json:"commit_ref"`
	SystemdServices *[]SystemdServiceGet `json:"systemd_services,omitempty"`
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

	var res ResPipelineGet

	p, err := repository.FindPipeline(name)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
		}
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	status := p.GetStatus()
	commitRef := p.GetCommitRef()

	systemdServices := &[]SystemdServiceGet{}
	res = ResPipelineGet{
		Name:            name,
		Status:          string(status),
		CommitRef:       commitRef,
		SystemdServices: systemdServices,
	}

	services, err := p.GetStatusSystemdServices()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	for _, s := range services {
		*res.SystemdServices = append(*res.SystemdServices, SystemdServiceGet{Name: name, Status: s.Status})
	}

	return c.JSONPretty(http.StatusOK, res, "	")
}
