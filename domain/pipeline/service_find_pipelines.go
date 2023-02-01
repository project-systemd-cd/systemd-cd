package pipeline

import "systemd-cd/domain/logger"

// FindPipelines implements IPipelineService
func (s pipelineService) FindPipelines() (metadatas []PipelineMetadata, err error) {
	logger.Logger().Debug("START - Find pipelines")
	defer func() {
		if err == nil {
			for i, pm := range metadatas {
				logger.Logger().Debugf("> pipeline[%d].Name = %v", i, pm.Name)
				logger.Logger().Debugf("> pipeline[%d].PathLocalRepository = %v", i, pm.PathLocalRepository)
			}
			logger.Logger().Debug("END   - Find pipelines")
		} else {
			logger.Logger().Error("FAILED - Find pipelines")
			logger.Logger().Error(err)
		}
	}()

	metadatas, err = s.repo.FindPipelines()
	return metadatas, err
}
