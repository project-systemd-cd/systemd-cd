package git_command

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"

	"gopkg.in/src-d/go-git.v4/config"
)

func (*GitCommand) SetRemoteUrl(workingDir git.Path, remoteName string, url string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "remoteName", Value: remoteName}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	err = r.DeleteRemote(remoteName)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
