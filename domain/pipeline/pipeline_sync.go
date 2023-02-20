package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/git"
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
			p.Status = StatusFailed
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
	updateExists, targetCommitId, targetTagName, err := p.updateExists()
	if err != nil {
		return err
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
	if oldStatus != StatusFailed {
		// Stop systemd service before backup
		var systemdServices []systemd.UnitService
		systemdServices, err = p.getSystemdServices()
		for _, s := range systemdServices {
			logger.Logger().Infof("Stop systemd unit service \"%v\"", s.Name)
			err = s.Disable(true)
			if err != nil {
				return err
			}
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

	if targetCommitId != nil {
		// Checkout
		err = p.RepositoryLocal.Checkout(*targetCommitId)
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
		_, err = p.RepositoryLocal.Pull()
		isFastForward := err == nil || !errorss.Is(err, git.ErrNonFastForwardUpdate)
		if err != nil && isFastForward {
			return err
		}
		if !isFastForward {
			_, err = p.RepositoryLocal.Reset(git.OptionReset{Mode: git.HardReset}, p.ManifestMerged.GitTargetBranch)
			if err != nil {
				return err
			}
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

	// Register jobs
	pipelineId := UUID()
	jobTest, err := p.newJobTest(pipelineId, targetTagName)
	if err != nil {
		return err
	}
	jobBuild, err := p.newJobBuild(pipelineId, targetTagName)
	if err != nil {
		return err
	}
	jobInstall, err := p.newJobInstall(pipelineId, targetTagName)
	if err != nil {
		return err
	}

	// Run jobs
	for _, job := range []*jobInstance{jobTest, jobBuild, jobInstall} {
		if job != nil {
			if err == nil {
				err2 := job.Run(p.service.repo)
				if err2 != nil {
					err = err2
				}
			} else {
				// if job failed, cancel reaming jobs
				err2 := job.Cancel(p.service.repo)
				if err2 != nil {
					err = err2
				}
			}
		}
	}
	if err != nil {
		return err
	}

	// Execute over systemd
	services, err := p.getSystemdServices()
	if err != nil {
		return err
	}
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
