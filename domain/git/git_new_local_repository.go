package git

const DefaultRemoteName = "origin"

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git gitService) NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error) {
	// Open git dir if exists
	exists, err := git.command.IsGitDirectory(path)
	if err != nil {
		return
	}
	if exists {
		// Git repository already exist
		var ref string
		ref, err = git.command.RefCommitId(path)
		if err != nil {
			return
		}
		// Check remote url
		var s string
		s, err = git.command.GetRemoteUrl(path, DefaultRemoteName)
		if err != nil {
			return
		}
		if s != remoteUrl {
			// if remote url is different, set remote url
			err = git.command.SetRemoteUrl(path, DefaultRemoteName, remoteUrl)
			if err != nil {
				return
			}
		}
		repo = &RepositoryLocal{
			git:          &git,
			RemoteUrl:    remoteUrl,
			TargetBranch: branch,
			RefCommitId:  ref,
			Path:         path,
		}

		return false, repo, nil
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
	repo = &RepositoryLocal{
		git:          &git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}

	return true, repo, nil
}
