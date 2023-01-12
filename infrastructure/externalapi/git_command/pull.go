package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func (g *GitCommand) Pull(workingDir git.Path, force bool) (refCommitId string, err error) {
	logger.Logger().Tracef("Called:\n\tworkingDir: %v\n\tforce: %v", workingDir, force)

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	w, err := r.Worktree()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	err = w.Pull(&gitcommand.PullOptions{Force: force})
	if err != nil && err != gitcommand.NoErrAlreadyUpToDate {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	r2, err := r.Head()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}

	refCommitId = r2.Hash().String()
	logger.Logger().Tracef("Finished:\n\trefCommitId: %v", refCommitId)
	return refCommitId, nil
}
