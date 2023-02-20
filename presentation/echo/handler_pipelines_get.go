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

	_, err = CheckJWT(c)
	if err != nil {
		err = c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
		return err
	}

	var res ResPipelinesGet = nil

	pp, err := service.FindPipelines()
	if err != nil {
		return err
	}
	for _, p := range pp {
		name := p.GetName()

		var ss []SystemdServiceGet = nil
		var systemdServices *[]SystemdServiceGet = &ss
		resItem := ResPipelineGet{
			Name:              name,
			GitRemoteUrl:      p.GetGitRemoteUrl(),
			GitTargetBranch:   p.GetGitTargetBranch(),
			GitTargetTagRegex: p.GetGitTargetTagRegex(),
			Status:            string(p.GetStatus()),
			AutoSyncEnabled:   p.AutoSyncEnabled(),
			CommitRef:         p.GetCommitRef(),
			SystemdServices:   systemdServices,
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

	if len(res) == 0 {
		err = c.NoContent(http.StatusNoContent)
		return err
	}
	err = c.JSONPretty(http.StatusOK, res, "	")
	return err
}
