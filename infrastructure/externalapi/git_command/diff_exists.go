package git_command

import (
	"systemd-cd/domain/git"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) DiffExists(workingDir git.Path, branch string) (exists bool, err error) {
	r, err := open(workingDir)
	if err != nil {
		return
	}
	headRef, err := r.Head()
	if err != nil {
		return
	}
	headCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		return
	}
	revHash, err := r.ResolveRevision(plumbing.Revision("origin/" + branch))
	if err != nil {
		return
	}
	revCommit, err := r.CommitObject(*revHash)
	if err != nil {
		return
	}

	exists = headCommit.Hash.String() != revCommit.Hash.String()
	return exists, nil
}
