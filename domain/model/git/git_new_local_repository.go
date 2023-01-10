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
