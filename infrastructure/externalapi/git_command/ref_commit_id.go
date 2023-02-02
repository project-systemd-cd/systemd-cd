package git_command

import (
	"systemd-cd/domain/git"
)

func (g *GitCommand) RefCommitId(workingDir git.Path) (string, error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	r2, err := r.Head()
	if err != nil {
		return "", err
	}

	refCommitId := r2.Hash().String()
	return refCommitId, nil
}
