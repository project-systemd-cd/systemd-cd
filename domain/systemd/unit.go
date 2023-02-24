package systemd

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

type IUnitService interface {
	GetName() string
	GetUnitFilePath() string

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
	StatusStopped  Status = "stopped"
	StatusRunning  Status = "running"
	StatusFailed   Status = "failed"
	StatusNotFound Status = "not found"
)

var (
	// check implements
	_ IUnitService = unitService{}
)

type (
	unitService struct {
		systemctl             Systemctl
		Name                  string
		unitFile              UnitFileService
		Path                  string
		EnvironmentFileValues map[string]string
	}
)

func (u unitService) GetName() string {
	return u.Name
}
func (u unitService) GetUnitFilePath() string {
	return u.Path
}

func (u unitService) Enable(startNow bool) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Enable systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debugf("< startNow = %v", startNow)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debug("END   - Enable systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Enable systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	err = u.systemctl.Enable(u.Name, startNow)
	return err
}

func (u unitService) Disable(stopNow bool) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Disable systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debugf("< stopNow = %v", stopNow)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debug("END   - Disable systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Disable systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	err = u.systemctl.Disable(u.Name, stopNow)
	return err
}

func (u unitService) Start() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Start systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debug("END   - Start systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Start systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	return u.systemctl.Start(u.Name)
}

func (u unitService) Stop() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Stop systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debug("END   - Stop systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Stop systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	err = u.systemctl.Stop(u.Name)
	return err
}

func (u unitService) Restart() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Restart systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debug("END   - Restart systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Restart systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	err = u.systemctl.Restart(u.Name)
	return err
}

func (u unitService) GetStatus() (s Status, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get status of systemd unit service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debugf("> status = %v", s)
			logger.Logger().Debug("END   - Get status of systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Get status of systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if u.Path == "" {
		return StatusNotFound, &errors.ErrNotFound{Object: "unit file", IdName: "path", Id: u.Name}
	}
	s, err = u.systemctl.Status(u.Name)
	return s, err
}
