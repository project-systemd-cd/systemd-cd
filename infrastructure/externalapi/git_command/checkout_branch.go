package git_command

import (
	"fmt"
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) CheckoutBranch(workingDir git.Path, branch string) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	b := plumbing.ReferenceName(branch)
	if !b.IsBranch() {
		return fmt.Errorf("invalid branch name '%s'", branch)
	}
	err = w.Checkout(&gitcommand.CheckoutOptions{Branch: b})
	if err != nil {
		return err
	}

	return nil
}
