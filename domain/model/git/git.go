package git

type IService interface {
	// Open local git repository.
	// If local git repository does not exist, execute clone.
	NewLocalRepository(path Path, remoteUrl string, branch string) (cloned bool, repo *RepositoryLocal, err error)
}

func NewService(git GitCommand) IService {
	return &gitService{command: git}
}

type gitService struct {
	command GitCommand
}
