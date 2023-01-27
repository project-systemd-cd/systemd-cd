package pipeline

import (
	"errors"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/systemd"
)

func (p *pipeline) Init() (err error) {
	logger.Logger().Tracef("Called:\n\tpipeline: %+v", p)

	defer func() {
		if err != nil {
			p.Status = StatusError
		}
	}()

	p.Status = StatusSyncing

	// Get manifest and merge local manifest
	m, err := p.loadManifest()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}
	p.ManifestMerged = m

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

	// Install
	service, err := p.install()
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Execute over systemd
	err = service.Start()
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
		err = errors.New("failed to ")
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	p.Status = StatusSynced
	logger.Logger().Tracef("Finished: \n\tpipeline.Status", StatusSynced)
	return nil
}
