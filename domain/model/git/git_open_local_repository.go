package git

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
