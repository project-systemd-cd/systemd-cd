package git

import "systemd-cd/domain/logger"

type IService interface {
	// Open local git repository.
	// If local git repository does not exist, execute clone.
	NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error)
}

func NewService(git GitCommand) IService {
	logger.Logger().Debug("START - Instantiate git service")
	defer func() {
		logger.Logger().Debug("END   - Instantiate git service")
	}()

	return &gitService{command: git}
}

type gitService struct {
	command GitCommand
}
