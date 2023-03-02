package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
)

func (p *pipeline) Init() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Initialize pipeline")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Initialize pipeline")
		} else {
			logger.Logger().Error("FAILED - Initialize pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	defer func() {
		if err != nil {
			p.Status = StatusFailed
		}
	}()

	p.Status = StatusSyncing

	_, targetCommitId, targetTagName, err := p.updateExists()
	if err != nil {
		return err
	}
	if targetCommitId != nil {
		// Checkout
		err = p.RepositoryLocal.Checkout(*targetCommitId)
		if err != nil {
			return err
		}
	} else {
		// Checkout branch
		err = p.RepositoryLocal.CheckoutBranch("refs/heads/" + p.ManifestLocal.GitTargetBranch)
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
			_, err = p.RepositoryLocal.Reset(git.OptionReset{Mode: git.HardReset}, p.ManifestLocal.GitTargetBranch)
			if err != nil {
				return err
			}
		}
	}

	// Get manifest and merge local manifest
	m, err := p.getRemoteManifest()
	if err != nil {
		return err
	}
	mm, err := m.merge(p.RepositoryLocal.RemoteUrl, p.ManifestLocal)
	if err != nil {
		return err
	}
	p.ManifestMerged = mm

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

	p.Status = StatusSynced
	return nil
}
