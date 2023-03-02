package runner

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
)

func (s *service) NewPipeline(m pipeline.ServiceManifestLocal, o OptionAddPipeline) (Pipeline, error) {
	p, err := s.FindPipeline(m.Name)
	if err == nil {
		// Remove pipeline if exist
		err := s.repository.RemovePipeline(m.Name)
		if err != nil {
			return Pipeline{}, err
		}
	}

	p1, err := s.pipelineService.NewPipeline(m)
	if err != nil {
		return Pipeline{}, err
	}

	p2 := Pipeline{p1, o.AutoSync}
	_, err = s.repository.AddPipeline(p2)

	if p.AutoSyncEnabled() != o.AutoSync {
		if o.AutoSync {
			logger.Logger().Infof("Pipeline \"%s\" auto sync enabled", m.Name)
		} else {
			logger.Logger().Infof("Pipeline \"%s\" auto sync disabled", m.Name)
		}
	}

	return p2, err
}
