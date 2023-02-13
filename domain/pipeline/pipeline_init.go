package pipeline

import (
	"fmt"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"time"
)

func (p *pipeline) Init() (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Initialize pipeline")
	logger.Logger().Infof("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Initialize pipeline")
		} else {
			logger.Logger().Error("FAILED - Initialize pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	defer func() {
		if err != nil {
			p.Status = StatusFailed
		}
	}()

	p.Status = StatusSyncing

	// Get manifest and merge local manifest
	logger.Logger().Info("Get manifest in git repository and merge to local manifest")
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
		err = us.Enable(true)
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

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
