package runner

func (s *service) RemovePipeline(name string) error {
	p, err := s.repository.FindPipeline(name)
	if err != nil {
		return err
	}
	err = p.Uninstall()
	if err != nil {
		return err
	}
	return s.repository.RemovePipeline(name)
}
