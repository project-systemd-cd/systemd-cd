package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Find git commit hash by regex of git tag")
	logger.Logger().Debugf("< force = %v", force)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> commitId = %v", refCommitId)
			logger.Logger().Debug("END   - Find git commit hash by regex of git tag")
		} else {
			logger.Logger().Error("FAILED - Find git commit hash by regex of git tag")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		return "", err
	}

	r.RefCommitId = refCommitId

	return refCommitId, nil
}
