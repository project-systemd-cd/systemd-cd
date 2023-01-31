package git_command

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) CheckoutHash(workingDir git.Path, hash string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "hash", Value: hash}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	err = w.Checkout(&gitcommand.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
