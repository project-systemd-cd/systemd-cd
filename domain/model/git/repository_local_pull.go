package git

import "systemd-cd/domain/model/logger"

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	logger.Logger().Trace("Called:\n\tforce: %v", force)

	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	r.RefCommitId = refCommitId

	logger.Logger().Tracef("Finished:\n\trefCommitId: %v", refCommitId)
	return
}
