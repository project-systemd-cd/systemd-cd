package git_command

import (
	"systemd-cd/domain/git"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (*GitCommand) GetCommitMessage(workingDir git.Path, hash string) (msg string, err error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	c, err := r.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		return "", err
	}
	return c.Message, nil
}
