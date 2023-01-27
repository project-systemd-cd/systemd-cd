package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) DiffExists(workingDir git.Path, branch string) (exists bool, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "branch", Value: branch}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	headRef, err := r.Head()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	headCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	revHash, err := r.ResolveRevision(plumbing.Revision("origin/" + branch))
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	revCommit, err := r.CommitObject(*revHash)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	exists = headCommit.Hash.String() != revCommit.Hash.String()
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "exists", Value: exists}}))
	return exists, nil
}
