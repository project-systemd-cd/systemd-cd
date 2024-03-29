package git

import "systemd-cd/domain/logger"

const DefaultRemoteName = "origin"

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git gitService) NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate git local repository")
	logger.Logger().Debugf("< path = %v", path)
	logger.Logger().Debugf("< remoteUrl = %v", remoteUrl)
	logger.Logger().Debugf("< branch = %v", branch)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> cloned = %v", cloned)
			logger.Logger().Debugf("> commitId = %v", repo.RefCommitId)
			logger.Logger().Tracef("> repositoryLocal = %+v", repo)
			logger.Logger().Debug("END   - Instantiate git local repository")
		} else {
			logger.Logger().Error("FAILED - Instantiate git local repository")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

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
		var currentRemote string
		currentRemote, err = git.command.GetRemoteUrl(path, DefaultRemoteName)
		if err != nil {
			return
		}
		if currentRemote != remoteUrl {
			// if remote url is different, set remote url
			err = git.command.SetRemoteUrl(path, DefaultRemoteName, remoteUrl)
			if err != nil {
				return
			}
		}
		// Check branch
		var currentBranch string
		currentBranch, err = git.command.RefBranchName(path)
		if currentBranch != branch {
			// if branch is different, checkout branch
			git.command.CheckoutBranch(path, branch)
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
