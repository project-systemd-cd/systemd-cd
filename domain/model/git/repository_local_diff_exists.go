package git

import "systemd-cd/domain/model/logger"

func (r *RepositoryLocal) DiffExists(executeFetch bool) (exists bool, err error) {
	logger.Logger().Trace("Called:\n\texecuteFetch: %v", executeFetch)
	if executeFetch {
		err = r.fetch()
		if err != nil {
			return
		}
	}

	exists, err = r.git.command.DiffExists(r.Path, r.TargetBranch)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	logger.Logger().Tracef("Finished:\n\texists: %v", exists)
	return
}
