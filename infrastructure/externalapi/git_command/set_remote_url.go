package git_command

import (
	"systemd-cd/domain/git"

	"gopkg.in/src-d/go-git.v4/config"
)

func (*GitCommand) SetRemoteUrl(workingDir git.Path, remoteName string, url string) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	err = r.DeleteRemote(remoteName)
	if err != nil {
		return err
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})
	if err != nil {
		return err
	}

	return nil
}
