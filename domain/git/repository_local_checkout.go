package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Checkout(hash string) (err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "hash", Value: hash}}))

	err = r.git.command.CheckoutHash(r.Path, hash)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	r.RefCommitId = hash

	logger.Logger().Trace("Finished")
	return
}
