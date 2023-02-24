package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) getSystemdServices() (systemdServices []systemd.IUnitService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get systemd unit services on pipeline")
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			for _, s := range systemdServices {
				logger.Logger().Debugf("> service.Name = %s", s.GetName())
				logger.Logger().Tracef("> service = %+v", s)
			}
			logger.Logger().Debug("END   - Get systemd unit services on pipeline")
		} else {
			logger.Logger().Error("FAILED - Get systemd unit services on pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	for _, service := range p.ManifestMerged.SystemdServiceOptions {
		var s systemd.IUnitService
		s, err2 := p.service.Systemd.GetService(service.Name)
		if err2 == nil {
			systemdServices = append(systemdServices, s)
		} else {
			err = err2
			systemdServices = append(systemdServices, s)
		}
	}

	return systemdServices, err
}
