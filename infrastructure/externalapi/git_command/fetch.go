package git_command

import (
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Fetch(workingDir git.Path) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	err = r.Fetch(&gitcommand.FetchOptions{Tags: gitcommand.AllTags})
	if err != nil && err != gitcommand.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}
