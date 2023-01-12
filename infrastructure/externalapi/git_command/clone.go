package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) Clone(path git.Path, remoteUrl string, targetBranch string, recursive bool) error {
	logger.Logger().Tracef("Called:\n\tpath: %v\n\tremoteUrl: %v\n\ttargetBranch: %v\n\trecursive: %v", path, remoteUrl, targetBranch, recursive)

	_, err := gitcommand.PlainClone(string(path), false, &gitcommand.CloneOptions{
		URL:               remoteUrl,
		ReferenceName:     plumbing.NewBranchReferenceName(targetBranch),
		RecurseSubmodules: gitcommand.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return err
}
