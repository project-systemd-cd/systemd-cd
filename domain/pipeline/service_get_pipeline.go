package pipeline

import (
	"systemd-cd/domain/logger"
)

// GetPipeline implements iPipelineService
func (s pipelineService) GetPipeline(name string) (p IPipeline, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate pipeline from repository")
	logger.Logger().Tracef("* pipelineService = %+v", s)
	logger.Logger().Debugf("< name = %v", name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate pipeline from repository")
		} else {
			logger.Logger().Error("FAILED - Instantiate pipeline from repository")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	//* NOTE: Receiver must not a pointer

	// Find pipeline from repository
	mp, err := s.repo.FindPipelineByName(name)
	if err != nil {
		return &pipeline{}, err
	}
	// Change directories temporary
	//* NOTE: this asignment only effect once to execute this func
	// TODO: move files and restart service
	s.PathSrcDir = mp.Service.PathSrcDir
	s.PathBinDir = mp.Service.PathBinDir
	s.PathEtcDir = mp.Service.PathEtcDir
	s.PathOptDir = mp.Service.PathOptDir
	s.PathSystemdUnitFileDir = mp.Service.PathSystemdUnitFileDir
	s.PathSystemdUnitEnvFileDir = mp.Service.PathSystemdUnitEnvFileDir
	s.PathBackupDir = mp.Service.PathBackupDir

	// Define pipeline
	p1 := &pipeline{
		ManifestLocal:   mp.ManifestLocal,
		RepositoryLocal: nil,
		Status:          StatusOutOfSync,
		service:         &s,
	}

	logger.Logger().Infof("Pipeline \"%v\" loaded", name)

	// Open local repository
	// if local repository not exists, clone remote repository
	_, p1.RepositoryLocal, err = s.Git.NewLocalRepository(mp.PathLocalRepository, mp.ManifestLocal.GitRemoteUrl, mp.ManifestLocal.GitTargetBranch)
	if err != nil {
		return &pipeline{}, err
	}

	// Get manifest and merge local manifest
	m2, err := p1.getRemoteManifest()
	if err != nil {
		return &pipeline{}, err
	}
	mm2, err := m2.merge(p1.RepositoryLocal.RemoteUrl, p1.ManifestLocal)
	if err != nil {
		return &pipeline{}, err
	}
	p1.ManifestMerged = mm2

	// Check updates
	updateExists, _, _, err := p1.updateExists()
	if err != nil {
		return &pipeline{}, err
	}
	if !updateExists {
		p1.Status = StatusSynced
	}

	p = p1
	return p, err
}
