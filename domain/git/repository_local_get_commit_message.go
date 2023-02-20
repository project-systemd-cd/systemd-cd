package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) GetCommitMessage(hash string) (msg string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get git commit message")
	logger.Logger().Debugf("< commitId = %v", hash)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> message = %v", msg)
			logger.Logger().Debug("END   - Get git commit message")
		} else {
			logger.Logger().Error("FAILED - Get git commit message")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	msg, err = r.git.command.GetCommitMessage(r.Path, hash)
	if err != nil {
		return "", err
	}

	return msg, err
}
