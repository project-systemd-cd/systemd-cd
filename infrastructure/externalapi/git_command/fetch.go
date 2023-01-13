package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Fetch(workingDir git.Path) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	err = r.Fetch(&gitcommand.FetchOptions{})
	if err == gitcommand.NoErrAlreadyUpToDate {
		logger.Logger().Trace("Finished")
		return nil
	}
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return err
}
