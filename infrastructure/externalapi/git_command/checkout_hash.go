package git_command

import (
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) CheckoutHash(workingDir git.Path, hash string) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Checkout(&gitcommand.CheckoutOptions{Hash: plumbing.NewHash(hash)})
	if err != nil {
		return err
	}

	return nil
}
