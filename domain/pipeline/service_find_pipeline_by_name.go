package pipeline

import "systemd-cd/domain/logger"

// FindPipelineByName implements IPipelineService
func (s pipelineService) FindPipelineByName(name string) (m PipelineMetadata, err error) {
	logger.Logger().Debug("START - Find pipeline by name")
	logger.Logger().Tracef("* pipelineService = %+v", s)
	logger.Logger().Debugf("< name = %v", name)
	defer func() {
		if err == nil {
			logger.Logger().Debugf("> pipelineMetadata.PathLocalRepository = %v", m.PathLocalRepository)
			logger.Logger().Tracef("> pipelineMetadata = %+v", m)
			logger.Logger().Debug("END   - Find pipeline by name")
		} else {
			logger.Logger().Error("FAILED - Find pipeline by name")
			logger.Logger().Error(err)
		}
	}()

	m, err = s.repo.FindPipelineByName(name)
	return m, err
}
