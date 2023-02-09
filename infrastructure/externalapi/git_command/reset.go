package git_command

import (
	"systemd-cd/domain/git"

	ggit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) Reset(workingDir git.Path, o git.OptionReset, target string) (refCommitId string, err error) {
	r, err := open(workingDir)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	revHash, err := r.ResolveRevision(plumbing.Revision(target))
	if err != nil {
		return
	}
	err = w.Reset(&ggit.ResetOptions{
		Commit: *revHash,
		Mode:   o.Mode,
	})
	if err != nil {
		return
	}

	r2, err := r.Head()
	if err != nil {
		return
	}
	refCommitId = r2.Hash().String()

	return refCommitId, nil
}
