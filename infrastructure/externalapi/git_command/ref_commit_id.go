package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) RefCommitId(workingDir git.Path) (string, error) {
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

	refCommitId := r2.Hash().String()
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "refCommitId", Value: refCommitId}}))
	return refCommitId, nil
}
