package runner_gitops

import (
	"reflect"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/runner"
	"time"
)

func (s *service) Start(option runner.Option) (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Start gitops pipeline runner")
	logger.Logger().Debugf("< option = %+v", option)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Start gitops pipeline runner")
		} else {
			logger.Logger().Error("FAILED - Start gitops pipeline runner")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	c := make(chan error)
	go func() {
		var manifests []pipeline.ServiceManifestLocal
		manifests, err = s.loadManifests()
		if err != nil {
			c <- err
		}
		option1 := option
		option1.RemovePipelineManifestFileNotSpecified = false
		c <- s.runner.Start(manifests, option1)
	}()

	for s.runner.IsLoading() {
		// Wait runner initialization
		time.Sleep(time.Second)
	}

	var prevManifests []pipeline.ServiceManifestLocal
	pipelines, err := s.runner.FindPipelines()
	for _, p := range pipelines {
		prevManifests = append(prevManifests, p.GetManifestLocal())
	}
	for {
		logger.Logger().Debug("-----------------------------------------------------------")
		logger.Logger().Debug("START - GitOps sync")
		logger.Logger().Debug("-----------------------------------------------------------")

		select {
		case err = <-c:
			return err
		default:
		}

		// Load manifest
		_, err = s.repository.Pull()
		if err != nil {
			return err
		}
		var manifests []pipeline.ServiceManifestLocal
		manifests, err = s.loadManifests()
		if err != nil {
			return err
		}

		for _, m := range manifests {
			found := false
			for _, prev := range prevManifests {
				if prev.Name == m.Name {
					if !reflect.DeepEqual(m, prev) {
						logger.Logger().Infof("Manifest file updated (name: %s)", m.Name)
						// Manifest file updated
						// Update pipeline
						_, err = s.runner.NewPipeline(m, runner.OptionAddPipeline{AutoSync: true})
						if err != nil {
							return err
						}
					}
					found = true
					break
				}
			}
			if !found {
				logger.Logger().Infof("New manifest file (name: %s)", m.Name)
				// New manifest file
				// Register pipeline
				_, err = s.runner.NewPipeline(m, runner.OptionAddPipeline{AutoSync: true})
				if err != nil {
					return err
				}
			}
		}
		// Disable auto sync from pipeline manifest file deleted
		for _, prev := range prevManifests {
			found := false
			for _, m := range manifests {
				if m.Name == prev.Name {
					found = true
					break
				}
			}
			if !found {
				logger.Logger().Infof("Manifest file removed (name: %s)", prev.Name)
				if option.RemovePipelineManifestFileNotSpecified {
					// Remove pipeline
					err = s.runner.RemovePipeline(prev.Name)
					if err != nil {
						return err
					}
				} else {
					// Disable pipeline auto sync
					var p runner.Pipeline
					p, err = s.runner.FindPipeline(prev.Name)
					if err != nil {
						return err
					}
					p.DisableAutoSync()
					logger.Logger().Infof("Pipeline \"%s\" auto sync disabled", prev.Name)
				}
			}
		}

		prevManifests = manifests

		logger.Logger().Debug("-----------------------------------------------------------")
		logger.Logger().Debug("END - GitOps sync")
		logger.Logger().Debug("-----------------------------------------------------------")

		time.Sleep(option.PollingInterval)
	}
}
