package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) Sync() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Sync pipeline")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Sync pipeline")
		} else {
			logger.Logger().Error("FAILED - Sync pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	defer func() {
		if err != nil {
			p.Status = StatusFailed
		}
	}()

	if p.Status == StatusSyncing {
		logger.Logger().Debugf("Skip to sync pipeline \"%v\", because state is syncing", p.ManifestLocal.Name)
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
	if !updateExists && (p.Status == StatusSynced || p.Status == StatusFailed) {
		// Already synced
		if p.ManifestMerged.GitTagRegex == nil {
			logger.Logger().Debugf("Pipeline \"%v\" has no updates (branch: %v)", p.ManifestMerged.Name, p.ManifestMerged.GitTargetBranch)
		} else {
			logger.Logger().Debugf("Pipeline \"%v\" has no updates (tag: %v)", p.ManifestMerged.Name, *p.ManifestMerged.GitTagRegex)
		}
		return nil
	}

	// Update exists
	logger.Logger().Infof("Pipeline \"%v\" is syncing", p.ManifestMerged.Name)
	oldStatus := p.Status
	oldCommitId := p.GetCommitRef()
	p.Status = StatusSyncing

	// Backup
	if oldStatus != StatusFailed {
		// Stop systemd service before backup
		var systemdServices []systemd.IUnitService
		systemdServices, err = p.getSystemdServices()
		for _, s := range systemdServices {
			logger.Logger().Debugf("Stop systemd unit service \"%v\"", s.GetName())
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
				if err != nil {
					// Restore from backup
					err = p.restoreBackup(restoreBackupOptions{&oldCommitId})
					if err != nil {
						return err
					}
					var services []systemd.IUnitService
					services, err = p.getSystemdServices()
					if err != nil {
						return err
					}
					for _, s := range services {
						// Restart systemd service
						err = s.Enable(true)
						if err != nil {
							return err
						}
					}
					// Checkout old commit id
					err = p.RepositoryLocal.Checkout(oldCommitId)
					if err != nil {
						return err
					}
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

	p.Status = StatusSynced
	return nil
}
