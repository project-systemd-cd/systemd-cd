package runner

import (
	"errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"
)

var pipelines []pipeline.IPipeline

func (s *runnerService) Start(manifests *[]pipeline.ServiceManifestLocal) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Info("START - Pipeline runner")
	if manifests != nil {
		for i, sml := range *manifests {
			logger.Logger().Debugf("> localManifest[%d].Name = %v", i, sml.Name)
			logger.Logger().Debugf("> localManifest[%d].GitRemoteUrl = %v", i, sml.GitRemoteUrl)
			logger.Logger().Debugf("> localManifest[%d].GitTargetBranch = %v", i, sml.GitTargetBranch)
			if sml.GitTagRegex == nil {
				logger.Logger().Debugf("> localManifest[%d].GitTagRegex = %v", i, nil)
			} else {
				logger.Logger().Debugf("> localManifest[%d].GitTagRegex = %v", i, *sml.GitTagRegex)
			}
			logger.Logger().Tracef("> localManifest[%d] = %+v", i, sml)
		}
	}
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Pipeline runner")
		} else {
			logger.Logger().Error("FAILED - Pipeline runner")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if manifests == nil {
		err = errors.New("\"manifests\" must not nil pointer")
		if err != nil {
			return err
		}
	}

	logger.Logger().Info("Get pipelines from repository")

	foundPipelines := []string{}
	metadatas, err := s.pipelineService.FindPipelines()
	if err != nil {
		return err
	}
	for _, m := range metadatas {
		manifestFileSpecified := false
		for _, sml := range *manifests {
			if sml.Name == m.Name {
				var i pipeline.IPipeline
				i, err = s.pipelineService.NewPipeline(sml)
				if err != nil {
					return err
				}
				err = i.Sync()
				if err != nil {
					return err
				}
				pipelines = append(pipelines, i)
				foundPipelines = append(foundPipelines, m.Name)
				manifestFileSpecified = true
				break
			}
		}
		if !manifestFileSpecified {
			var ip pipeline.IPipeline
			ip, err = s.pipelineService.NewPipeline(m.ManifestLocal)
			if err != nil {
				return err
			}
			// TODO: Don't sync pipelines that don't have specified manifest files periodically
			pipelines = append(pipelines, ip)
		}
	}
	for _, m := range *manifests {
		found := false
		for _, foundPipeline := range foundPipelines {
			if m.Name == foundPipeline {
				found = true
				break
			}
		}
		if !found {
			logger.Logger().Infof("Initialize pipeline \"%v\"", m.Name)
			var pipeline pipeline.IPipeline
			pipeline, err = s.pipelineService.NewPipeline(m)
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
