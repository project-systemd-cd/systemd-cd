package git

import "systemd-cd/domain/model/logger"

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git gitService) NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error) {
	logger.Logger().Tracef("Called:\n\tpath: %v\n\tremoteUrl: %v\n\tbranch: %v", path, remoteUrl, branch)

	// Open git dir if exists
	exists, err := git.command.IsGitDirectory(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	if exists {
		// Git repository already exist
		var ref string
		ref, err = git.command.RefCommitId(path)
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return
		}
		repo = &RepositoryLocal{
			git:          &git,
			RemoteUrl:    remoteUrl,
			TargetBranch: branch,
			RefCommitId:  ref,
			Path:         path,
		}
		logger.Logger().Tracef("Finished:\n\tRepositoryLocal: %v", *repo)
		return false, repo, nil
	}

	// Clone
	err = git.command.Clone(path, remoteUrl, branch, true)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}

	// Get ref
	ref, err := git.command.RefCommitId(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	repo = &RepositoryLocal{
		git:          &git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}
	logger.Logger().Tracef("Finished:\n\tRepositoryLocal: %v", *repo)
	return true, repo, nil
}
