package pipeline

import (
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
)

// NewPipeline implements iPipelineService
func (s pipelineService) NewPipeline(m ServiceManifestLocal) (iPipeline, error) {
	logger.Logger().Tracef("Called:\n\tmanifest: %+v", m)

	// Define pipeline
	p := &pipeline{
		ManifestLocal:   m,
		RepositoryLocal: nil,
		Status:          StatusOutOfSync,
		service:         &s,
	}

	// Open local repository
	var err error
	var cloned bool
	// if local repository not exists, clone remote repository
	cloned, p.RepositoryLocal, err = s.Git.NewLocalRepository(git.Path(s.PathSrcDir+m.Name), m.GitRemoteUrl, m.GitTargetBranch)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return &pipeline{}, err
	}

	if cloned {
		err = p.Init()
	} else {
		err = p.Sync()
	}
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return p, err
	}

	logger.Logger().Tracef("Finished:\n\tpipeline: %+v", *p)
	return p, nil
}
