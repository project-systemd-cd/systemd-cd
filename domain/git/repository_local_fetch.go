package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Fetch() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Git fetch")
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Git fetch")
		} else {
			logger.Logger().Error("FAILED - Git fetch")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = r.git.command.Fetch(r.Path)
	return err
}
