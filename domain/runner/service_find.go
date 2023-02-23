package runner

func (s *service) FindPipeline(name string) (Pipeline, error) {
	return s.repository.FindPipeline(name)
}

func (s *service) FindPipelines() ([]Pipeline, error) {
	return s.repository.FindPipelines()
}
