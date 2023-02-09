package git

import (
	"systemd-cd/domain/logger"
)

func (r *RepositoryLocal) Reset(o OptionReset, branch string) (refCommitId string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Git reset")
	logger.Logger().Debugf("< option = %+v", o)
	logger.Logger().Debugf("< branch = %s", branch)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> commitId = %v", refCommitId)
			logger.Logger().Debug("END   - Git reset")
		} else {
			logger.Logger().Error("FAILED - Git reset")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	refCommitId, err = r.git.command.Reset(r.Path, o, "origin/"+branch)
	if err != nil {
		return "", err
	}

	r.RefCommitId = refCommitId

	return refCommitId, nil
}
