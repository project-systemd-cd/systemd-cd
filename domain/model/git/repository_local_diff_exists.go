package git

import "systemd-cd/domain/model/logger"

func (r *RepositoryLocal) DiffExists(executeFetch bool) (exists bool, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "executeFetch", Value: executeFetch}}))
	if executeFetch {
		err = r.fetch()
		if err != nil {
			return
		}
	}

	exists, err = r.git.command.DiffExists(r.Path, r.TargetBranch)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "exists", Value: exists}}))
	return
}
