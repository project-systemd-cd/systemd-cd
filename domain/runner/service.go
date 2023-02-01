package runner

import (
	"systemd-cd/domain/pipeline"
	"time"
)

type IRunnerService interface {
	Start(*[]pipeline.ServiceManifestLocal) error
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

func NewService(p pipeline.IPipelineService, o Option) (IRunnerService, error) {
	err := o.validate()
	if err != nil {
		return &runnerService{}, err
	}

	return &runnerService{p, o}, nil
}

type runnerService struct {
	pipelineService pipeline.IPipelineService

	option Option
}
