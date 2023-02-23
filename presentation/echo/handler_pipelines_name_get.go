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
	QueryParamPipelineJobsGet
}

type ResPipelineGet struct {
	Name              string               `json:"name"`
	GitRemoteUrl      string               `json:"git_remote_url"`
	GitTargetBranch   string               `json:"git_target_branch"`
	GitTargetTagRegex *string              `json:"git_target_tag_regex,omitempty"`
	Status            string               `json:"status"`
	AutoSyncEnabled   bool                 `json:"auto_sync"`
	CommitRef         string               `json:"commit_ref"`
	SystemdServices   *[]SystemdServiceGet `json:"systemd_services,omitempty"`
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

	logger.Logger().Trace("-----------------------------------------------------------")
	logger.Logger().Trace("START - GET /pipelines/:name")
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
			logger.Logger().Tracef("END    - GET /pipelines/:name %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - GET /pipelines/:name")
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

	query := QueryParamPipelineGet{}
	if err = c.Bind(&query); err != nil {
		return err
	}
	query.Embed = c.QueryParams()["embed"]

	p, err := service.FindPipeline(name)
	if err != nil {
		var ErrNotFound *errors.ErrNotFound
		if errorss.As(err, &ErrNotFound) {
			err = c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
			return err
		}
		return err
	}

	var ss []SystemdServiceGet = nil
	var systemdServices *[]SystemdServiceGet = &ss
	res := ResPipelineGet{
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

		var jobs [][]pipeline.Job
		jobs, err = p.GetJobs(query2)
		if err != nil {
			return err
		}

		err = c.JSONPretty(http.StatusOK, ResPipelineGetJobsEmbed{res, jobs}, "	")
		return err
	} else {
		err = c.JSONPretty(http.StatusOK, res, "	")
		return err
	}
}
