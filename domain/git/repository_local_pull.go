package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	logger.Logger().Debug("START - Find git commit hash by regex of git tag")
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Find git commit hash by regex of git tag")
		} else {
			logger.Logger().Error("FAILED - Find git commit hash by regex of git tag")
			logger.Logger().Error(err)
		}
	}()

	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		return "", err
	}

	r.RefCommitId = refCommitId

	return refCommitId, nil
}
