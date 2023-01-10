package git

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git *Git) NewLocalRepository(path Path, remoteUrl string, branch string) (repo *RepositoryLocal, err error) {
	// Open git dir if exists
	_, err = git.command.Status(path)
	if err != ErrRepositoryNotExists {
		var ref string
		ref, err = git.command.RefCommitId(path)
		if err != nil {
			return
		}
		return &RepositoryLocal{
			git:          git,
			RemoteUrl:    remoteUrl,
			TargetBranch: branch,
			RefCommitId:  ref,
			Path:         path,
		}, nil
	}

	// Clone
	err = git.command.Clone(path, remoteUrl, branch, true)
	if err != nil {
		return
	}

	// Get ref
	ref, err := git.command.RefCommitId(path)
	if err != nil {
		return
	}
	return &RepositoryLocal{
		git:          git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}, err
}

// Open local git repository.
func (git *Git) OpenLocalRepository(path Path) (repo *RepositoryLocal, err error) {
	// Open git dir if exists
	_, err = git.command.Status(path)
	if err != nil {
		return nil, err
	}

	remoteUrl, err := git.command.GetRemoteUrl(path, "origin")
	if err != nil {
		return nil, err
	}
	branch, err := git.command.RefBranchName(path)
	if err != nil {
		return nil, err
	}
	ref, err := git.command.RefCommitId(path)
	if err != nil {
		return nil, err
	}

	return &RepositoryLocal{
		git:          git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}, nil
}

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		return
	}
	r.RefCommitId = refCommitId
	return
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
