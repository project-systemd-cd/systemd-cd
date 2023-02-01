package git_command

import (
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) Clone(path git.Path, remoteUrl string, targetBranch string, recursive bool) error {
	_, err := gitcommand.PlainClone(string(path), false, &gitcommand.CloneOptions{
		URL:               remoteUrl,
		ReferenceName:     plumbing.NewBranchReferenceName(targetBranch),
		RecurseSubmodules: gitcommand.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}

	return nil
}
