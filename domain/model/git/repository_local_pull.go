package git

import "systemd-cd/domain/model/logger"

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "force", Value: force}}))

	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	r.RefCommitId = refCommitId

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "refCommitId", Value: refCommitId}}))
	return
}
