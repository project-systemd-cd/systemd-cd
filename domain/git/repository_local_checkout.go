package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Checkout(hash string) (err error) {
	logger.Logger().Debug("START - Git checkout to commit")
	logger.Logger().Debugf("< commitId = %v", hash)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Git checkout to commit")
		} else {
			logger.Logger().Error("FAILED - Git checkout to commit")
			logger.Logger().Error(err)
		}
	}()

	err = r.git.command.CheckoutHash(r.Path, hash)
	if err != nil {
		return
	}

	r.RefCommitId = hash

	return
}
