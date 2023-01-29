package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) Sync() (err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	// Get manifest and merge local manifest
	m, err := p.loadManifest()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	p.ManifestMerged = m

	// Check update
	if outOfSync, err := p.RepositoryLocal.DiffExists(true); err != nil {
		return err
	} else if !outOfSync {
		// Already synced
		p.Status = StatusSynced
		logger.Logger().Tracef("Finished: Pipeline \"%s\" already up to date", p.ManifestMerged.Name)
		return nil
	}

	// Update exists
	if p.Status == StatusSyncing {
		logger.Logger().Debugf("Pipeline \"%s\" is syncing", p.ManifestMerged.Name)
		return nil
	}
	oldStatus := p.Status
	p.Status = StatusSyncing

	// Pull
	_, err = p.RepositoryLocal.Pull(false)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Get manifest and merge local manifest
	m2, err := p.loadManifest()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	p.ManifestMerged = m2

	// Test
	err = p.test()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Build
	err = p.build()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Backup
	if oldStatus != StatusError {
		// TODO: stop systemd service before backup
		err = p.backupInstalled()
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	// Install
	services, err := p.install()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	logger.Logger().Debugf("Debug:\n\tservices: %v", services)

	for _, s := range services {
		// Execute over systemd
		err = s.Restart()
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}

		// Get status of systemd service
		s, err := s.GetStatus()
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
		if s != systemd.StatusRunning {
			// If failed to execute over systemd, restore from backup
			err = p.restoreBackup(restoreBackupOptions{})
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}

	p.Status = StatusSynced
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: p.Status}}))
	return nil
}
