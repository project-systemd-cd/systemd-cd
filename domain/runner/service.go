package runner

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"
)

type IRunnerService interface {
	Start(*[]pipeline.ServiceManifestLocal) error

	// Resgister pipeline.
	// If pipeline with same name already exists, replace it.
	NewPipeline(pipeline.ServiceManifestLocal, OptionAddPipeline) (Pipeline, error)

	FindPipeline(name string) (Pipeline, error)
	FindPipelines() ([]Pipeline, error)
}

type OptionAddPipeline struct {
	AutoSync bool
}

func NewService(p pipeline.IPipelineService, o Option) (service IRunnerService, err error) {
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

	service = &runnerService{p, o, inmemoryRepository()}
	return service, err
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

type runnerService struct {
	pipelineService pipeline.IPipelineService

	option     Option
	repository iRepositoryInmemory
}
