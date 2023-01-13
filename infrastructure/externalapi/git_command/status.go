package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) IsGitDirectory(workingDir git.Path) (bool, error) {
	logger.Logger().Tracef("Called:\n\tworkingDir: %v", workingDir)

	_, err := open(workingDir)
	if err == git.ErrRepositoryNotExists {
		logger.Logger().Tracef("Finished:\n\tIsGitDirectory: %v", false)
		return false, nil
	}
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return false, err
	}
	logger.Logger().Tracef("Finished:\n\tIsGitDirectory: %v", true)
	return true, nil
}
