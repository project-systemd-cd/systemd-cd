package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) CheckoutBranch(branch string) (err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "branch", Value: branch}}))

	err = r.git.command.CheckoutBranch(r.Path, branch)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	s, err := r.git.command.RefCommitId(r.Path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	r.RefCommitId = s

	logger.Logger().Trace("Finished")
	return
}
