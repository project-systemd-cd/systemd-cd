package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) getSystemdServices() (systemdServices []systemd.UnitService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get systemd unit services on pipeline")
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			for _, s := range systemdServices {
				logger.Logger().Debugf("> service.Name = %s", s.Name)
				logger.Logger().Tracef("> service = %+v", s)
			}
			logger.Logger().Debug("END   - Get systemd unit services on pipeline")
		} else {
			logger.Logger().Error("FAILED - Get systemd unit services on pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	for _, service := range p.ManifestMerged.SystemdOptions {
		var s systemd.UnitService
		s, err = p.service.Systemd.GetService(service.Name)
		if err != nil {
			return nil, err
		}
		systemdServices = append(systemdServices, s)
	}

	return systemdServices, nil
}