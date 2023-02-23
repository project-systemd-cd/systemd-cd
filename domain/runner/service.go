package runner

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"
)

type IRunnerService interface {
	Start([]pipeline.ServiceManifestLocal, Option) error

	// Resgister pipeline.
	// If pipeline with same name already exists, replace it.
	NewPipeline(pipeline.ServiceManifestLocal, OptionAddPipeline) (Pipeline, error)

	FindPipeline(name string) (Pipeline, error)
	FindPipelines() ([]Pipeline, error)
	RemovePipeline(name string) error
}

type Option struct {
	PollingInterval                        time.Duration
	RemovePipelineManifestFileNotSpecified bool
}

type OptionAddPipeline struct {
	AutoSync bool
}

func NewService(ps pipeline.IPipelineService) (s IRunnerService) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate runner service")
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		logger.Logger().Debug("END   - Instantiate runner service")
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	return &service{ps, inmemoryRepository()}
}

type service struct {
	pipelineService pipeline.IPipelineService

	repository iRepositoryInmemory
}
