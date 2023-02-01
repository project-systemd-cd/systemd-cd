package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) HeadIsLatesetOfBranch(branch string) (isLatest bool, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Judge git head is latest of branch")
	logger.Logger().Debugf("< branch = %v", branch)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Judge git head is latest of branch")
			logger.Logger().Debugf("> isLatest = %v", isLatest)
		} else {
			logger.Logger().Error("FAILED - Judge git head is latest of branch")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	exists, err := r.git.command.DiffExists(r.Path, branch)
	if err != nil {
		return false, err
	}
	isLatest = !exists

	return isLatest, nil
}
