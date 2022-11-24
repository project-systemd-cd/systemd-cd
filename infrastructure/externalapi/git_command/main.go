package git_command

import (
	"systemd-cd/domain/model/git"

	gitcommand "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func New() git.GitCommand {
	var g git.GitCommand = &GitCommand{}
	return g
}

// implements "systemd-cd/domain/model/git".GitCommand
type GitCommand struct{}

func (g *GitCommand) Clone(path git.Path, remoteUrl string, targetBranch string, recursive bool) error {
	_, err := gitcommand.PlainClone(string(path), false, &gitcommand.CloneOptions{
		URL:               remoteUrl,
		ReferenceName:     plumbing.NewBranchReferenceName(targetBranch),
		RecurseSubmodules: gitcommand.DefaultSubmoduleRecursionDepth,
	})
	return err
}

func (g *GitCommand) Fetch(workingDir git.Path) error {
	r, err := open(workingDir)
	if err != nil {
		return err
	}
	err = r.Fetch(&gitcommand.FetchOptions{})
	if err == gitcommand.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (g *GitCommand) DiffExists(workingDir git.Path, to string) (exists bool, err error) {
	r, err := open(workingDir)
	if err != nil {
		return
	}
	headRef, err := r.Head()
	if err != nil {
		return
	}
	headCommit, err := r.CommitObject(headRef.Hash())
	if err != nil {
		return
	}
	revHash, err := r.ResolveRevision(plumbing.Revision("origin/" + to))
	if err != nil {
		return
	}
	revCommit, err := r.CommitObject(*revHash)
	if err != nil {
		return
	}
	return headCommit.Hash.String() != revCommit.Hash.String(), nil
}

func (g *GitCommand) Pull(workingDir git.Path, force bool) (refCommitId string, err error) {
	r, err := open(workingDir)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	err = w.Pull(&gitcommand.PullOptions{Force: force})
	if err != nil && err != gitcommand.NoErrAlreadyUpToDate {
		return
	}
	r2, err := r.Head()
	if err != nil {
		return
	}
	return r2.Hash().String(), nil
}

func (g *GitCommand) Status(workingDir git.Path) (string, error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	s, err := w.Status()
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

func (g *GitCommand) Ref(workingDir git.Path) (string, error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}
	r2, err := r.Reference(plumbing.HEAD, false)
	if err != nil {
		return "", err
	}
	return r2.Hash().String(), nil
}

func open(dir git.Path) (r *gitcommand.Repository, err error) {
	r, err = gitcommand.PlainOpen(string(dir))
	if err == gitcommand.ErrRepositoryNotExists {
		err = git.ErrRepositoryNotExists
	}
	return
}
