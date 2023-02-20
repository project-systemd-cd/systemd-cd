package git

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
)

var (
	ErrRepositoryNotExists  = errors.New("repository does not exist")
	ErrNonFastForwardUpdate = git.ErrNonFastForwardUpdate
)

type GitCommand interface {
	Clone(path Path, remoteUrl string, targetBranch string, recursive bool) error
	Fetch(workingDir Path) error
	DiffExists(workingDir Path, to string) (exists bool, err error)
	Pull(workingDir Path, force bool) (refCommitId string, err error)
	Reset(workingDir Path, o OptionReset, target string) (refCommitId string, err error)
	IsGitDirectory(workingDir Path) (bool, error)
	RefCommitId(workingDir Path) (string, error)
	RefBranchName(workingDir Path) (string, error)
	FindHashByTagRegex(workingDir Path, regex string) (hash string, err error)
	CheckoutBranch(workingDir Path, branch string) error
	CheckoutHash(workingDir Path, hash string) error
	GetCommitMessage(workingDir Path, hash string) (string, error)
	GetRemoteUrl(workingDir Path, remoteName string) (string, error)
	SetRemoteUrl(workingDir Path, remoteName string, url string) error
}

type OptionReset struct {
	Mode ResetMode
}

type ResetMode = git.ResetMode

const (
	MixedReset = git.MixedReset
	HardReset  = git.HardReset
	MergeReset = git.MergeReset
	SoftReset  = git.SoftReset
)
