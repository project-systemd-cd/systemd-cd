package git_command

import (
	"errors"
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) GetRemoteUrl(workingDir git.Path, remoteName string) (string, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "workingDir", Value: workingDir}, {Name: "remoteName", Value: remoteName}}))

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	r2, err := r.Remote(remoteName)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	s := r2.Config().URLs
	if len(s) == 0 {
		err := errors.New("invalid remote url list length")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "GitRemoteUrl", Value: s[0]}}))
	return s[0], nil
}
