package systemd

import "systemd-cd/domain/logger"

type Unit interface {
	Enable(startNow bool) error
	Disable(stopNow bool) error
	Start() error
	Stop() error
	Restart() error
	GetStatus() (Status, error)
}

type (
	Status string
)

const (
	// Systemd service status
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
	StatusFailed  Status = "failed"
)

var (
	// check implements
	_ Unit = UnitService{}
)

type (
	UnitService struct {
		systemctl             Systemctl
		Name                  string
		unitFile              UnitFileService
		Path                  string
		EnvironmentFileValues map[string]string
	}
)

func (u UnitService) Enable(startNow bool) (err error) {
	logger.Logger().Debug("START - Enable systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debugf("< startNow = %v", startNow)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Enable systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Enable systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	err = u.systemctl.Enable(u.Name, startNow)
	return err
}

func (u UnitService) Disable(stopNow bool) (err error) {
	logger.Logger().Debug("START - Disable systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debugf("< stopNow = %v", stopNow)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Disable systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Disable systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	err = u.systemctl.Disable(u.Name, stopNow)
	return err
}

func (u UnitService) Start() (err error) {
	logger.Logger().Debug("START - Start systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Start systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Start systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	return u.systemctl.Start(u.Name)
}

func (u UnitService) Stop() (err error) {
	logger.Logger().Debug("START - Stop systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Stop systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Stop systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	err = u.systemctl.Stop(u.Name)
	return err
}

func (u UnitService) Restart() (err error) {
	logger.Logger().Debug("START - Restart systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Restart systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Restart systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	err = u.systemctl.Restart(u.Name)
	return err
}

func (u UnitService) GetStatus() (s Status, err error) {
	logger.Logger().Debug("START - Get status of systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debugf("> status = %v", s)
			logger.Logger().Debug("END   - Get status of systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Get status of systemd unit service")
			logger.Logger().Error(err)
		}
	}()

	s, err = u.systemctl.Status(u.Name)
	return s, err
}
