package runner

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"
)

type IRunnerService interface {
	Start(*[]pipeline.ServiceManifestLocal) error

	FindPipeline(name string) (pipeline.IPipeline, error)
	FindPipelines() ([]pipeline.IPipeline, error)
}

type Option struct {
	PollingInterval time.Duration
}

func (o Option) validate() error {
	if o.PollingInterval < 3*time.Minute {
		// return errors.New("polling interval must be at least 3 minutes")
	}
	return nil
}

func NewService(p pipeline.IPipelineService, repo IRepositoryInmemory, o Option) (service IRunnerService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate pipeline service")
	logger.Logger().Debugf("< option = %+v", o)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate pipeline service")
		} else {
			logger.Logger().Error("FAILED - Instantiate pipeline service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = o.validate()
	if err != nil {
		return &runnerService{}, err
	}

	service = &runnerService{p, o, repo}
	return service, err
}

type runnerService struct {
	pipelineService pipeline.IPipelineService

	option     Option
	repository IRepositoryInmemory
}
