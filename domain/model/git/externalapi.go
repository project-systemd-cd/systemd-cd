package git

import (
	"errors"
)

var (
	ErrRepositoryNotExists = errors.New("repository does not exist")
)

type GitCommand interface {
	Clone(path Path, remoteUrl string, targetBranch string, recursive bool) error
	Fetch(workingDir Path) error
	DiffExists(workingDir Path, to string) (exists bool, err error)
	Pull(workingDir Path, force bool) (refCommitId string, err error)
	Status(workingDir Path) (string, error)
	Ref(workingDir Path) (string, error)
}
