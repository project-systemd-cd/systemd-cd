package pipeline

import "systemd-cd/domain/logger"

// FindPipelines implements IPipelineService
func (s pipelineService) FindPipelines() (metadatas []PipelineMetadata, err error) {
	logger.Logger().Debug("START - Find pipelines")
	logger.Logger().Tracef("* pipelineService = %+v", s)
	defer func() {
		if err == nil {
			for i, pm := range metadatas {
				logger.Logger().Debugf("> pipelineMetadata[%d].Name = %v", i, pm.Name)
				logger.Logger().Debugf("> pipelineMetadata[%d].PathLocalRepository = %v", i, pm.PathLocalRepository)
				logger.Logger().Tracef("> pipelineMetadata[%d] = %+v", i, pm)
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
