package pipeline

import (
	errorss "errors"
	"reflect"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
)

// NewPipeline implements iPipelineService
func (s pipelineService) NewPipeline(m ServiceManifestLocal) (p IPipeline, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate pipeline with repository data")
	logger.Logger().Tracef("* pipelineService = %+v", s)
	logger.Logger().Debugf("< manifestLocal.Name = %v", m.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate pipeline with repository data")
		} else {
			logger.Logger().Error("FAILED - Instantiate pipeline with repository data")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	//* NOTE: Receiver must not a pointer

	// TODO: validate local manifest

	// Find pipeline from repository
	manifestConfigured := false
	mp, err := s.repo.FindPipelineByName(m.Name)
	ErrNotFound := &errors.ErrNotFound{}
	notFound := errorss.As(err, &ErrNotFound)
	if err != nil && !notFound {
		return &pipeline{}, err
	}
	if notFound {
		mp.Name = m.Name
		mp.PathLocalRepository = git.Path(s.PathSrcDir + m.Name)
		mp.ManifestLocal = m
		mp.Service.PathSrcDir = s.PathSrcDir
		mp.Service.PathBinDir = s.PathBinDir
		mp.Service.PathEtcDir = s.PathEtcDir
		mp.Service.PathOptDir = s.PathOptDir
		mp.Service.PathSystemdUnitFileDir = s.PathSystemdUnitFileDir
		mp.Service.PathSystemdUnitEnvFileDir = s.PathSystemdUnitEnvFileDir
		mp.Service.PathBackupDir = s.PathBackupDir

		// Save to repository
		err = s.repo.SavePipeline(mp)
		if err != nil {
			return &pipeline{}, err
		}
	} else {
		if reflect.DeepEqual(m, mp.ManifestLocal) {
			logger.Logger().Infof("Pipeline \"%v\" loaded", m.Name)
		} else {
			logger.Logger().Infof("Pipeline \"%s\" configured", m.Name)
			manifestConfigured = true

			// if manifests are different, replace them
			mp.ManifestLocal = m

			// Save to repository
			err = s.repo.SavePipeline(mp)
			if err != nil {
				return &pipeline{}, err
			}
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
	}

	// Define pipeline
	p1 := &pipeline{
		ManifestLocal:   mp.ManifestLocal,
		RepositoryLocal: nil,
		Status:          StatusOutOfSync,
		service:         &s,
	}

	// Open local repository
	var cloned bool
	// if local repository not exists, clone remote repository
	cloned, p1.RepositoryLocal, err = s.Git.NewLocalRepository(mp.PathLocalRepository, mp.ManifestLocal.GitRemoteUrl, mp.ManifestLocal.GitTargetBranch)
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

	if manifestConfigured {
		err = p1.Uninstall()
		if err != nil {
			return &pipeline{}, err
		}
	}
	if notFound || cloned || manifestConfigured {
		err = p1.Init()
		if err != nil {
			err = p1.Uninstall()
			if err != nil {
				return &pipeline{}, err
			}
		}
	} else {
		// Check updates
		updateExists, _, _, err := p1.updateExists()
		if err != nil {
			return &pipeline{}, err
		}
		if !updateExists {
			p1.Status = StatusSynced
		}
	}

	p = p1
	return p, err
}
