package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) HeadIsLatesetOfBranch(branch string) (bool, error) {
	logger.Logger().Trace("Called")

	exists, err := r.git.command.DiffExists(r.Path, branch)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return false, err
	}
	isLatest := !exists

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "isLatest", Value: isLatest}}))
	return isLatest, nil
}
