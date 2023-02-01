package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) CheckoutBranch(branch string) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Git checkout to brach")
	logger.Logger().Debugf("< branch = %v", branch)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Git checkout to brach")
		} else {
			logger.Logger().Error("FAILED - Git checkout to brach")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = r.git.command.CheckoutBranch(r.Path, branch)
	if err != nil {
		return
	}

	s, err := r.git.command.RefCommitId(r.Path)
	if err != nil {
		return
	}

	r.RefCommitId = s

	return
}
