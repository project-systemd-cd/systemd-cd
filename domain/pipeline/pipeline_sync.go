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

	if p.Status == StatusSyncing {
		logger.Logger().Debugf("Pipeline \"%s\" is syncing", p.ManifestMerged.Name)
		return nil
	}

	// Get manifest and merge local manifest
	m, err := p.getRemoteManifest()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	mm, err := m.merge(p.RepositoryLocal.RemoteUrl, p.ManifestLocal)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	p.ManifestMerged = mm

	// Check update
	updateExists, err := p.GetUpdateExistence()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if !updateExists {
		// Already synced
		p.Status = StatusSynced
		logger.Logger().Tracef("Finished: Pipeline \"%s\" already up to date", p.ManifestMerged.Name)
		return nil
	}

	// Update exists
	oldStatus := p.Status
	oldCommitId := p.GetCommitRef()
	p.Status = StatusSyncing

	// Backup
	if oldStatus != StatusError {
		// TODO: stop systemd service before backup
		err = p.backupInstalled()
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	// Pull
	_, err = p.RepositoryLocal.Pull(false)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Get manifest and merge local manifest
	m2, err := p.getRemoteManifest()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	mm2, err := m2.merge(p.RepositoryLocal.RemoteUrl, p.ManifestLocal)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	p.ManifestMerged = mm2

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

	// Install
	services, err := p.install()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	logger.Logger().Debugf("Debug:\n\tservices: %v", services)

	// Execute over systemd
	if services != nil || len(services) != 0 {
		failedToExecuteOverSystemd := false
		for _, s := range services {
			// Execute over systemd
			err = s.Restart()
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}

			// Get status of systemd service
			status, err := s.GetStatus()
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
			for _, s := range services {
				// Restart systemd service
				err = s.Restart()
				if err != nil {
					logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
					return err
				}
			}
			// TODO: checkout old commit id
			// TODO: record commit id failed
		}
	}

	p.Status = StatusSynced
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: p.Status}}))
	return nil
}
