package runner

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"time"
)

func (s *service) Start(manifests []pipeline.ServiceManifestLocal, option Option) (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Pipeline runner")
	for i, sml := range manifests {
		logger.Logger().Debugf("< localManifest[%d].Name = %v", i, sml.Name)
		logger.Logger().Debugf("< localManifest[%d].GitRemoteUrl = %v", i, sml.GitRemoteUrl)
		logger.Logger().Debugf("< localManifest[%d].GitTargetBranch = %v", i, sml.GitTargetBranch)
		if sml.GitTagRegex == nil {
			logger.Logger().Debugf("< localManifest[%d].GitTagRegex = %v", i, nil)
		} else {
			logger.Logger().Debugf("< localManifest[%d].GitTagRegex = %v", i, *sml.GitTagRegex)
		}
		logger.Logger().Tracef("< localManifest[%d] = %+v", i, sml)
	}
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Pipeline runner")
		} else {
			logger.Logger().Error("FAILED - Pipeline runner")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	foundPipelines := []string{}
	metadatas, err := s.pipelineService.FindPipelines()
	if err != nil {
		return err
	}
	for _, savedManifest := range metadatas {
		manifestSpecified := false
		for _, m := range manifests {
			if m.Name == savedManifest.Name {
				var p Pipeline
				p, err = s.NewPipeline(m, OptionAddPipeline{AutoSync: true})
				if err != nil {
					return err
				}
				err = p.Sync()
				if err != nil {
					return err
				}
				foundPipelines = append(foundPipelines, m.Name)
				manifestSpecified = true
				break
			}
		}
		if !manifestSpecified {
			_, err = s.NewPipeline(
				savedManifest.ManifestLocal,
				OptionAddPipeline{AutoSync: false},
			)
			if err != nil {
				return err
			}
			if option.RemovePipelineManifestFileNotSpecified {
				err = s.RemovePipeline(savedManifest.ManifestLocal.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	for _, m := range manifests {
		found := false
		for _, foundPipeline := range foundPipelines {
			if m.Name == foundPipeline {
				found = true
				break
			}
		}
		if !found {
			logger.Logger().Infof("Initialize pipeline \"%v\"", m.Name)
			_, err = s.NewPipeline(m, OptionAddPipeline{AutoSync: true})
			if err != nil {
				return err
			}
		}
	}

	s.loading = false

	time.Sleep(option.PollingInterval)

	for {
		var pipelines []Pipeline
		pipelines, err = s.repository.FindPipelines()
		if err != nil {
			return err
		}
		for _, p := range pipelines {
			if p.AutoSyncEnabled() {
				err = p.Sync()
				if err != nil {
					return err
				}
			}
		}

		time.Sleep(option.PollingInterval)
	}
}
