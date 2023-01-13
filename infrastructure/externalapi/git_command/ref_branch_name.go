package git_command

import (
	"errors"
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) RefBranchName(workingDir git.Path) (string, error) {
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
	if !r2.Name().IsBranch() {
		return "", errors.New("ref `" + r2.String() + "` is not git branch")
	}

	logger.Logger().Tracef("Finished:\n\trefBranchName: %v", r2.String())
	return r2.String(), nil
}
