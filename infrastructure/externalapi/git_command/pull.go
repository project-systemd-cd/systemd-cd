package git_command

import (
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Pull(workingDir git.Path, force bool) (refCommitId string, err error) {
	r, err := open(workingDir)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	err = w.Pull(&gitcommand.PullOptions{Force: force})
	if err != nil && err != gitcommand.NoErrAlreadyUpToDate {
		return
	}
	r2, err := r.Head()
	if err != nil {
		return
	}

	refCommitId = r2.Hash().String()
	return refCommitId, nil
}
