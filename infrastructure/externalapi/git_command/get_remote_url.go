package git_command

import (
	"errors"
	"systemd-cd/domain/git"
)

func (g *GitCommand) GetRemoteUrl(workingDir git.Path, remoteName string) (string, error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	r2, err := r.Remote(remoteName)
	if err != nil {
		return "", err
	}
	s := r2.Config().URLs
	if len(s) == 0 {
		err := errors.New("invalid remote url list length")
		return "", err
	}

	return s[0], nil
}
