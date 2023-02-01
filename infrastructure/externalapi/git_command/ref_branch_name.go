package git_command

import (
	"fmt"
	"systemd-cd/domain/git"
)

func (g *GitCommand) RefBranchName(workingDir git.Path) (string, error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	r2, err := r.Head()
	if err != nil {
		return "", err
	}
	if !r2.Name().IsBranch() {
		return "", fmt.Errorf("ref '%s' is not git branch", r2.Name().String())
	}

	return r2.Name().Short(), nil
}
