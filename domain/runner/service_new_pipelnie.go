package runner

import "systemd-cd/domain/pipeline"

func (s *service) NewPipeline(m pipeline.ServiceManifestLocal, o OptionAddPipeline) (Pipeline, error) {
	if _, err := s.FindPipeline(m.Name); err == nil {
		// Remove pipeline if exist
		err := s.repository.RemovePipeline(m.Name)
		if err != nil {
			return Pipeline{}, err
		}
	}

	p, err := s.pipelineService.NewPipeline(m)
	if err != nil {
		return Pipeline{}, err
	}

	p2 := Pipeline{p, o.AutoSync}
	_, err = s.repository.AddPipeline(p2)

	return p2, err
}
