package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"time"
)

func (p *pipeline) Sync() (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Sync pipeline")
	logger.Logger().Infof("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Sync pipeline")
		} else {
			logger.Logger().Error("FAILED - Sync pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	if p.Status == StatusSyncing {
		logger.Logger().Infof("Skip to sync pipeline \"%v\", because state is syncing", p.ManifestLocal.Name)
		return nil
	}

	// Get manifest and merge local manifest
	{
		var m ServiceManifestRemote
		m, err = p.getRemoteManifest()
		if err != nil {
			return err
		}
		mm, err := m.merge(p.RepositoryLocal.RemoteUrl, p.ManifestLocal)
		if err != nil {
			return err
		}
		p.ManifestMerged = mm
	}

	// Check updates
	err = p.RepositoryLocal.Fetch()
	if err != nil {
		return err
	}
	updateExists := false
	var checkoutCommitId *string
	if p.ManifestMerged.GitTagRegex != nil {
		hash, err := p.RepositoryLocal.FindHashByTagRegex(*p.ManifestLocal.GitTagRegex)
		if err != nil {
			var ErrNotFound *errors.ErrNotFound
			if !errorss.As(err, &ErrNotFound) {
				return err
			}
		} else {
			if hash != p.RepositoryLocal.RefCommitId {
				updateExists = true
				checkoutCommitId = &hash
			}
		}
	} else {
		// Check update
		// Check update
		latest, err := p.RepositoryLocal.HeadIsLatesetOfBranch(p.ManifestMerged.GitTargetBranch)
		if err != nil {
			return err
		}
		if !latest {
			updateExists = true
		}
	}
	if !updateExists {
		// Already synced
		if p.ManifestMerged.GitTagRegex == nil {
			logger.Logger().Infof("Pipeline \"%v\" has no updates (branch: %v)", p.ManifestMerged.Name, p.ManifestMerged.GitTargetBranch)
		} else {
			logger.Logger().Infof("Pipeline \"%v\" has no updates (tag: %v)", p.ManifestMerged.Name, *p.ManifestMerged.GitTagRegex)
		}
		p.Status = StatusSynced
		return nil
	}

	// Update exists
	oldStatus := p.Status
	oldCommitId := p.GetCommitRef()
	p.Status = StatusSyncing

	// Backup
	if oldStatus != StatusError {
		for _, s := range p.ManifestMerged.SystemdOptions {
			// TODO: stop systemd service before backup
			logger.Logger().Infof("Stop systemd unit service \"%v\"", s.Name)
		}
		_, err = p.findBackupByCommitId(p.RepositoryLocal.RefCommitId)
		var ErrNotFound *errors.ErrNotFound
		notFound := errorss.As(err, &ErrNotFound)
		if notFound {
			err = p.backupInstalled()
			if err != nil {
				return err
			}
		} else {
			if err != nil {
				return err
			}
		}
	}

	if checkoutCommitId != nil {
		// Checkout
		err = p.RepositoryLocal.Checkout(*checkoutCommitId)
		if err != nil {
			return err
		}
	} else {
		// Checkout branch
		err = p.RepositoryLocal.CheckoutBranch("refs/heads/" + p.ManifestMerged.GitTargetBranch)
		if err != nil {
			return err
		}
		// Pull
		_, err = p.RepositoryLocal.Pull(false)
		if err != nil {
			return err
		}
	}

	// Get manifest and merge local manifest
	m2, err := p.getRemoteManifest()
	if err != nil {
		return err
	}
	mm2, err := m2.merge(p.RepositoryLocal.RemoteUrl, p.ManifestLocal)
	if err != nil {
		return err
	}
	p.ManifestMerged = mm2

	// Test
	err = p.test()
	if err != nil {
		return err
	}

	// Build
	err = p.build()
	if err != nil {
		return err
	}

	// Install
	services, err := p.install()
	if err != nil {
		return err
	}

	// Execute over systemd
	if services != nil || len(services) != 0 {
		failedToExecuteOverSystemd := false
		for _, s := range services {
			// Execute over systemd
			err = s.Restart()
			if err != nil {
				return err
			}

			time.Sleep(time.Second)

			// Get status of systemd service
			status, err := s.GetStatus()
			if err != nil {
				return err
			}
			if status != systemd.StatusRunning {
				failedToExecuteOverSystemd = true
				break
			}
		}
		if failedToExecuteOverSystemd {
			// Restore from backup
			err = p.restoreBackup(restoreBackupOptions{&oldCommitId})
			if err != nil {
				return err
			}
			for _, s := range services {
				// Restart systemd service
				err = s.Restart()
				if err != nil {
					return err
				}
			}
			// Checkout old commit id
			err = p.RepositoryLocal.Checkout(oldCommitId)
			if err != nil {
				return err
			}
			// TODO: record commit id failed
		}
	}

	p.Status = StatusSynced
	return nil
}
