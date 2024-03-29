package git

import (
	"errors"
	"systemd-cd/domain/logger"
)

func (r *RepositoryLocal) Pull() (refCommitId string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Git pull")
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		isNonFastForward := errors.Is(err, ErrNonFastForwardUpdate)
		if err == nil {
			logger.Logger().Debugf("> commitId = %v", refCommitId)
			logger.Logger().Debug("END   - Git pull")
		} else if isNonFastForward {
			logger.Logger().Debug("Skipped because non-fast-forward")
			logger.Logger().Debug("END   - Git pull")
		} else {
			logger.Logger().Error("FAILED - Git pull")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	refCommitId, err = r.git.command.Pull(r.Path, false)
	if err != nil {
		return "", err
	}

	r.RefCommitId = refCommitId

	return refCommitId, nil
}
