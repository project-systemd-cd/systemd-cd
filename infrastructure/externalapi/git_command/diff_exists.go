package git_command

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (g *GitCommand) DiffExists(workingDir git.Path, branch string) (exists bool, err error) {
	logger.Logger().Tracef("Called:\n\tworkingDir: %v\n\tbranch: %v", workingDir, branch)

	r, err := open(workingDir)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	headRef, err := r.Head()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	headCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	revHash, err := r.ResolveRevision(plumbing.Revision("origin/" + branch))
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	revCommit, err := r.CommitObject(*revHash)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}

	exists = headCommit.Hash.String() != revCommit.Hash.String()
	logger.Logger().Tracef("Finished:\n\texists: %v", exists)
	return exists, nil
}
