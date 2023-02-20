package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) GetCommitAuthor(hash string) (author string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get git commit author")
	logger.Logger().Debugf("< commitId = %v", hash)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> author = %v", author)
			logger.Logger().Debug("END   - Get git commit author")
		} else {
			logger.Logger().Error("FAILED - Get git commit author")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	author, err = r.git.command.GetCommitAuthor(r.Path, hash)
	if err != nil {
		return "", err
	}

	return author, err
}
