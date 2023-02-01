package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Checkout(hash string) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Git checkout to commit")
	logger.Logger().Debugf("< commitId = %v", hash)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Git checkout to commit")
		} else {
			logger.Logger().Error("FAILED - Git checkout to commit")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = r.git.command.CheckoutHash(r.Path, hash)
	if err != nil {
		return
	}

	r.RefCommitId = hash

	return
}
