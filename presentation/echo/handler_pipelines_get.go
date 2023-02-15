package echo

import (
	"net/http"
	"systemd-cd/domain/logger"

	"github.com/labstack/echo/v4"
)

type ResPipelinesGet []ResPipelineGet

func pipelinesGet(c echo.Context) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - GET /pipelines")
	logger.Logger().Tracef("< RemoteAddr = %s", c.Request().RemoteAddr)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> Status = %d", c.Response().Status)
			logger.Logger().Tracef("> ContentLength = %d", c.Response().Size)
			logger.Logger().Infof("END    - GET /pipelines %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	var res ResPipelinesGet = nil

	pp, err := repository.FindPipelines()
	if err != nil {
		return err
	}
	for _, p := range pp {
		name := p.GetName()
		status := p.GetStatus()
		commitRef := p.GetCommitRef()

		systemdServices := &[]SystemdServiceGet{}
		resItem := ResPipelineGet{
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
			*resItem.SystemdServices = append(*resItem.SystemdServices, SystemdServiceGet{Name: name, Status: s.Status})
		}

		res = append(res, resItem)
	}

	if res == nil {
		err = c.NoContent(http.StatusNoContent)
		return err
	}
	err = c.JSONPretty(http.StatusOK, res, "	")
	return err
}