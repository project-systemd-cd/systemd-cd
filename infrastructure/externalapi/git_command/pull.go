package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Pull(workingDir git.Path, force bool) (refCommitId string, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "force", Value: force}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	w, err := r.Worktree()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	err = w.Pull(&gitcommand.PullOptions{Force: force})
	if err != nil && err != gitcommand.NoErrAlreadyUpToDate {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	r2, err := r.Head()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	refCommitId = r2.Hash().String()
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "refCommitId", Value: refCommitId}}))
	return refCommitId, nil
}
