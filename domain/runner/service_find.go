package runner

import "systemd-cd/domain/pipeline"

func (s *runnerService) FindPipeline(name string) (pipeline.IPipeline, error) {
	return s.repository.FindPipeline(name)
}

func (s *runnerService) FindPipelines() ([]pipeline.IPipeline, error) {
	return s.repository.FindPipelines()
}
