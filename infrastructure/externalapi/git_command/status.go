package git_command

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
)

func (g *GitCommand) IsGitDirectory(workingDir git.Path) (bool, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}}))

	_, err := open(workingDir)
	if err == git.ErrRepositoryNotExists {
		logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "IsGitDirectory", Value: false}}))
		return false, nil
	}
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return false, err
	}
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "IsGitDirectory", Value: true}}))
	return true, nil
}
