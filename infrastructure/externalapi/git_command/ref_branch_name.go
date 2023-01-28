package git_command

import (
	"fmt"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
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
		return "", fmt.Errorf("ref '%s' is not git branch", r2.Name().String())
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "refBranchName", Value: r2.Name().Short()}}))
	return r2.Name().Short(), nil
}
