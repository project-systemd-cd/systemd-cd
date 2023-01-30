package pipeline

import (
	"fmt"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) Init() (err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	p.Status = StatusSyncing

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

	for _, s := range services {
		// Execute over systemd
		err = s.Start()
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
			err = fmt.Errorf("systemd service '%s' is not running", p.ManifestMerged.Name)
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	p.Status = StatusSynced
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: p.Status}}))
	return nil
}
