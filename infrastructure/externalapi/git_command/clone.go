package git_command

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) Clone(path git.Path, remoteUrl string, targetBranch string, recursive bool) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: path}, {Name: "remoteUrl", Value: remoteUrl}, {Name: "targetBranch", Value: targetBranch}, {Name: "recursive", Value: recursive}}))

	_, err := gitcommand.PlainClone(string(path), false, &gitcommand.CloneOptions{
		URL:               remoteUrl,
		ReferenceName:     plumbing.NewBranchReferenceName(targetBranch),
		RecurseSubmodules: gitcommand.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
