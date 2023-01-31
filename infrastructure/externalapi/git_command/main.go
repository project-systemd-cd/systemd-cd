package git_command

import (
	"systemd-cd/domain/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
)

func New() git.GitCommand {
	var g git.GitCommand = &GitCommand{}
	return g
}

// implements "systemd-cd/domain/git".GitCommand
type GitCommand struct{}

func open(dir git.Path) (r *gitcommand.Repository, err error) {
	r, err = gitcommand.PlainOpen(string(dir))
	if err == gitcommand.ErrRepositoryNotExists {
		err = git.ErrRepositoryNotExists
	}
	return
}
