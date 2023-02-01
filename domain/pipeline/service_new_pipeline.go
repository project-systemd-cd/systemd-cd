package pipeline

import (
	errorss "errors"
	"reflect"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/git"
)

// NewPipeline implements iPipelineService
func (s pipelineService) NewPipeline(m ServiceManifestLocal) (IPipeline, error) {
	//* NOTE: Receiver must not a pointer

	// Find pipeline from repository
	mp, err := s.repo.FindPipeline(m.Name)
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
		if !reflect.DeepEqual(m, mp.ManifestLocal) {
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
		s.PathSrcDir = mp.Service.PathSrcDir
		s.PathBinDir = mp.Service.PathBinDir
		s.PathEtcDir = mp.Service.PathEtcDir
		s.PathOptDir = mp.Service.PathOptDir
		s.PathSystemdUnitFileDir = mp.Service.PathSystemdUnitFileDir
		s.PathSystemdUnitEnvFileDir = mp.Service.PathSystemdUnitEnvFileDir
		s.PathBackupDir = mp.Service.PathBackupDir
	}

	// Define pipeline
	p := &pipeline{
		ManifestLocal:   mp.ManifestLocal,
		RepositoryLocal: nil,
		Status:          StatusOutOfSync,
		service:         &s,
	}

	// Open local repository
	var cloned bool
	// if local repository not exists, clone remote repository
	cloned, p.RepositoryLocal, err = s.Git.NewLocalRepository(mp.PathLocalRepository, mp.ManifestLocal.GitRemoteUrl, mp.ManifestLocal.GitTargetBranch)
	if err != nil {
		return &pipeline{}, err
	}

	if cloned {
		err = p.Init()
	} else {
		if p.Status != StatusSyncing {
			err = p.Sync()
		}
	}
	if err != nil {
		return p, err
	}

	return p, nil
}
