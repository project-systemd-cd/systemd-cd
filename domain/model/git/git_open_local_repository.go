package git

import "systemd-cd/domain/model/logger"

// Open local git repository.
func (git *Git) OpenLocalRepository(path Path) (repo *RepositoryLocal, err error) {
	logger.Logger().Tracef("Called:\n\tpath: %v", path)

	// Open git dir if exists
	exists, err := git.command.IsGitDirectory(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return
	}
	if !exists {
		// Git repository does not exist
		logger.Logger().Errorf("Error:\n\terr: %v", ErrRepositoryNotExists)
		return nil, ErrRepositoryNotExists
	}

	remoteUrl, err := git.command.GetRemoteUrl(path, "origin")
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return nil, err
	}
	branch, err := git.command.RefBranchName(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return nil, err
	}
	ref, err := git.command.RefCommitId(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return nil, err
	}

	repo = &RepositoryLocal{
		git:          git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}
	logger.Logger().Tracef("Finished:\n\tRepositoryLocal: %v", *repo)
	return repo, nil
}
