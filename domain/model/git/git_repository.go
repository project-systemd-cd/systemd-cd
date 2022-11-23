package git

import (
	"strings"
)

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git *Git) NewLocalRepository(path Path, remote RepositoryRemote, branch string) (repo *RepositoryLocal, err error) {
	// Open git dir if exists
	_, err = git.command.Status(path)
	if err != ErrRepositoryNotExists {
		var ref string
		ref, err = git.command.Ref(path)
		if err != nil {
			return
		}
		return &RepositoryLocal{
			git:              git,
			RepositoryRemote: remote,
			TargetBranch:     branch,
			RefCommitId:      ref,
			Path:             path,
		}, nil
	}

	// Generate remote url with access token
	remoteUrl := remote.RemoteUrl
	if remote.User != nil && remote.Token != nil && *remote.User != "" && *remote.Token != "" {
		remoteUrl = strings.Join(
			strings.Split(remote.RemoteUrl, "://"),
			"://"+*remote.User+":"+*remote.Token+"@",
		)
	}
	// Clone
	err = git.command.Clone(path, remoteUrl, branch, true)
	if err != nil {
		return
	}

	// Get ref
	ref, err := git.command.Ref(path)
	if err != nil {
		return
	}
	return &RepositoryLocal{
		git:              git,
		RepositoryRemote: remote,
		TargetBranch:     branch,
		RefCommitId:      ref,
		Path:             path,
	}, err
}

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	return r.git.command.Pull(r.Path, force)
}

func (r *RepositoryLocal) fetch() error {
	return r.git.command.Fetch(r.Path)
}

func (r *RepositoryLocal) DiffExists(executeFetch bool) (exists bool, err error) {
	if executeFetch {
		err = r.fetch()
		if err != nil {
			return
		}
	}
	return r.git.command.DiffExists(r.Path, r.TargetBranch)
}
