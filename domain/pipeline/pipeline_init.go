package pipeline

import (
	"fmt"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) Init() (err error) {
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

	for _, s := range services {
		// Execute over systemd
		err = s.Start()
		if err != nil {
			return err
		}

		// Get status of systemd service
		s, err := s.GetStatus()
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
