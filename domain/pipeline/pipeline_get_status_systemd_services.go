package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) GetStatusSystemdServices() (ss []SystemdServiceWithStatus, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get status of systemd unit services on pipeline")
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			for _, s := range ss {
				logger.Logger().Debugf("> service.Name = %s", s.Name)
				logger.Logger().Debugf("> service.Status = %v", s.Status)
				logger.Logger().Tracef("> service = %+v", s)
			}
			logger.Logger().Debug("END   - Get status of systemd unit services on pipeline")
		} else {
			logger.Logger().Error("FAILED - Get status of systemd unit services on pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	systemdServices, err := p.getSystemdServices()
	var ErrNotFound *errors.ErrNotFound
	if err != nil && !errorss.As(err, &ErrNotFound) {
		return nil, err
	}

	for _, s := range systemdServices {
		var status systemd.Status
		if s.Path == "" {
			status = systemd.StatusNotFound
		} else {
			status, err = s.GetStatus()
			if err != nil {
				return nil, err
			}
		}
		ss = append(ss, SystemdServiceWithStatus{s, status})
	}

	return ss, nil
}
