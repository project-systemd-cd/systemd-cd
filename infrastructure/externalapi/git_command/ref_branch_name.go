package git_command

import (
	"errors"
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) RefBranchName(workingDir git.Path) (string, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	r2, err := r.Head()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	if !r2.Name().IsBranch() {
		return "", errors.New("ref `" + r2.String() + "` is not git branch")
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "refBranchName", Value: r2.Name().Short()}}))
	return r2.Name().Short(), nil
}
