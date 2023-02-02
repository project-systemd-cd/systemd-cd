package git_command

import (
	"systemd-cd/domain/git"
)

func (g *GitCommand) IsGitDirectory(workingDir git.Path) (bool, error) {
	_, err := open(workingDir)
	if err == git.ErrRepositoryNotExists {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
