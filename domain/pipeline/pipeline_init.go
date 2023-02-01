package pipeline

import (
	"fmt"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) Init() (err error) {
	logger.Logger().Debug("START - Initialize pipeline")
	logger.Logger().Debugf("< pipeline.Name = %v", p.ManifestMerged.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Initialize pipeline")
		} else {
			logger.Logger().Error("FAILED - Initialize pipeline")
			logger.Logger().Error(err)
		}
	}()

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	p.Status = StatusSyncing

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

	for _, us := range services {
		// Execute over systemd
		err = us.Start()
		if err != nil {
			return err
		}

		// Get status of systemd service
		var s systemd.Status
		s, err = us.GetStatus()
		if err != nil {
			return err
		}
		if s != systemd.StatusRunning {
			err = fmt.Errorf("systemd service '%s' is not running", p.ManifestMerged.Name)
			return err
		}
	}

	p.Status = StatusSynced
	return nil
}
