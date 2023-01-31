package git_command

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Fetch(workingDir git.Path) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	err = r.Fetch(&gitcommand.FetchOptions{Tags: gitcommand.AllTags})
	if err == gitcommand.NoErrAlreadyUpToDate {
		logger.Logger().Trace("Finished")
		return nil
	}
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
