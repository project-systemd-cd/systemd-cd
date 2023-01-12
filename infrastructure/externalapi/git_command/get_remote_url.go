package git_command

import (
	"errors"
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

func (g *GitCommand) GetRemoteUrl(workingDir git.Path, remoteName string) (string, error) {
	logger.Logger().Tracef("Called:\n\tworkingDir: %v\n\tremoteName: %v", workingDir, remoteName)

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}
	r2, err := r.Remote(remoteName)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}
	s := r2.Config().URLs
	if len(s) == 0 {
		err := errors.New("invalid remote url list length")
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}

	logger.Logger().Tracef("Finished:\n\tGitRemoteUrl: %v", s[0])
	return s[0], nil
}
