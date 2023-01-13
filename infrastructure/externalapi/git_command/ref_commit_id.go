package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) RefCommitId(workingDir git.Path) (string, error) {
	logger.Logger().Tracef("Called:\n\tworkingDir: %v", workingDir)

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}
	r2, err := r.Head()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}

	refCommitId := r2.Hash().String()
	logger.Logger().Tracef("Finished\n\trefCommitId: %v", refCommitId)
	return refCommitId, nil
}
