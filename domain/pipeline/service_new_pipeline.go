package pipeline

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
)

// NewPipeline implements iPipelineService
func (s pipelineService) NewPipeline(m ServiceManifestLocal) (iPipeline, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: m}}))

	// TODO: find pipeline from repository.
	// TODO: if manifest updated, update pipeline in repository.

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
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return &pipeline{}, err
	}

	if cloned {
		err = p.Init()
	} else {
		err = p.Sync()
	}
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return p, err
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *p}}))
	return p, nil
}
