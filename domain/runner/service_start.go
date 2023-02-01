package runner

import (
	"errors"
	"systemd-cd/domain/pipeline"
	"time"
)

var pipelines []pipeline.IPipeline

func (s *runnerService) Start(manifests *[]pipeline.ServiceManifestLocal) error {
	if manifests == nil {
		err := errors.New("\"manifests\" must not nil pointer")
		if err != nil {
			return err
		}
	}

	foundPipelines := []string{}
	metadatas, err := s.pipelineService.FindPipelines()
	if err != nil {
		return err
	}
	for _, m := range metadatas {
		ip, err := s.pipelineService.NewPipeline(m.ManifestLocal)
		if err != nil {
			return err
		}

		for _, sml := range *manifests {
			if sml.Name == m.Name {
				pipelines = append(pipelines, ip)
				foundPipelines = append(foundPipelines, m.Name)
			}
		}
	}
	for _, m := range *manifests {
		found := false
		for _, foundPipeline := range foundPipelines {
			if m.Name == foundPipeline {
				found = true
			}
		}
		if !found {
			pipeline, err := s.pipelineService.NewPipeline(m)
			if err != nil {
				return err
			}
			err = pipeline.Sync()
			if err != nil {
				return err
			}
			pipelines = append(pipelines, pipeline)
		}
	}

	time.Sleep(s.option.PollingInterval)

	for {
		for _, p := range pipelines {
			err = p.Sync()
			if err != nil {
				return err
			}
		}

		time.Sleep(s.option.PollingInterval)
	}
}
