package runner_gitops

import (
	"strings"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/runner"
	"systemd-cd/domain/unix"
)

type IRunnerService interface {
	loadManifests() ([]pipeline.ServiceManifestLocal, error)

	Start(Option) error
}

type Option = runner.Option

func NewService(rs runner.IRunnerService, gs git.IService, dir Directory, r Repository) (s IRunnerService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate gitops runner service")
	logger.Logger().Debugf("> dir.Src = %s", dir.Src)
	logger.Logger().Debugf("> repository = %+v", r)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate gitops runner service")
		} else {
			logger.Logger().Error("FAILED - Instantiate gitops runner service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// Clone
	err = unix.MkdirIfNotExist(dir.Src)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(dir.Src, "/") {
		dir.Src += "/"
	}
	repoName := func() string {
		s := strings.Split(r.RemoteUrl, "/")
		s2 := strings.Split(s[len(s)-1], ".")
		return s2[0]
	}()
	_, rl, err := gs.NewLocalRepository(git.Path(dir.Src+repoName), r.RemoteUrl, r.Branch)
	if err != nil {
		return nil, err
	}

	return &service{rs, rl}, nil
}

type Directory struct {
	Src string
}

type Repository struct {
	RemoteUrl string
	Branch    string
}

type service struct {
	runner     runner.IRunnerService
	repository *git.RepositoryLocal
}
