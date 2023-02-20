package runner

func (s *runnerService) FindPipeline(name string) (Pipeline, error) {
	return s.repository.FindPipeline(name)
}

func (s *runnerService) FindPipelines() ([]Pipeline, error) {
	return s.repository.FindPipelines()
}
