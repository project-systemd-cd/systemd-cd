package git

import "systemd-cd/domain/logger"

const DefaultRemoteName = "origin"

// Open local git repository.
// If local git repository does not exist, execute clone.
func (git gitService) NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: path}, {Name: "remoteUrl", Value: remoteUrl}, {Name: "branch", Value: branch}}))

	// Open git dir if exists
	exists, err := git.command.IsGitDirectory(path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	if exists {
		// Git repository already exist
		var ref string
		ref, err = git.command.RefCommitId(path)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return
		}
		// Check remote url
		var s string
		s, err = git.command.GetRemoteUrl(path, DefaultRemoteName)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return
		}
		if s != remoteUrl {
			// if remote url is different, set remote url
			err = git.command.SetRemoteUrl(path, DefaultRemoteName, remoteUrl)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
		logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *repo}}))
		return false, repo, nil
	}

	// Clone
	err = git.command.Clone(path, remoteUrl, branch, true)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	// Get ref
	ref, err := git.command.RefCommitId(path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	repo = &RepositoryLocal{
		git:          &git,
		RemoteUrl:    remoteUrl,
		TargetBranch: branch,
		RefCommitId:  ref,
		Path:         path,
	}
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *repo}}))
	return true, repo, nil
}
