package pipeline

import (
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/systemd"
)

func (p *pipeline) Sync() (err error) {
	logger.Logger().Tracef("Called:\n\tpipeline: %+v", p)

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	// Get manifest and merge local manifest
	m, err := p.loadManifest()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
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
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Get manifest and merge local manifest
	m2, err := p.loadManifest()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}
	p.ManifestMerged = m2

	// Test
	err = p.test()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Build
	err = p.build()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Backup
	if oldStatus != StatusError {
		err = p.backupInstalled()
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	// Install
	service, err := p.install()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Execute over systemd
	err = service.Restart()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Get status of systemd service
	s, err := service.GetStatus()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}
	if s != systemd.StatusRunning {
		// If failed to execute over systemd, restore from backup
		err = p.restoreBackup(restoreBackupOptions{})
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	p.Status = StatusSynced
	logger.Logger().Tracef("Finished: \n\tpipeline.Status", StatusSynced)
	return nil
}
