package git_command

import (
	"fmt"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) Checkout(workingDir git.Path, branch string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "branch", Value: branch}}))

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
	b := plumbing.ReferenceName(branch)
	if !b.IsBranch() {
		return fmt.Errorf("invalid branch name '%s'", branch)
	}
	err = w.Checkout(&gitcommand.CheckoutOptions{Branch: b})
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
