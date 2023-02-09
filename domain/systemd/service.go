package systemd

import (
	"errors"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

var (
	ErrUnitFileNotManaged    = errors.New("unit file not managed by systemd-cd")
	ErrUnitEnvFileNotManaged = errors.New("unit env file not managed by systemd-cd")
)

type IService interface {
	// Generate unit-file.
	// If unit-file already exists, replace it.
	NewService(name string, uf UnitFileService, env map[string]string) (UnitService, error)
	GetService(name string) (UnitService, error)
	DeleteService(u UnitService) error

	loadUnitFileSerivce(path string) (u UnitFileService, isGeneratedBySystemdCd bool, err error)
	writeUnitFileService(u UnitFileService, path string) error

	loadEnvFile(path string) (e map[string]string, isGeneratedBySystemdCd bool, err error)
	writeEnvFile(e map[string]string, path string) error
}

func New(s Systemctl, unitFileDir string) (service IService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate systemd service (domain service, not systemd unit service)")
	logger.Logger().Debugf("< unitFileDir = %v", unitFileDir)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate systemd service (domain service, not systemd unit service)")
		} else {
			logger.Logger().Error("FAILED - Instantiate systemd service (domain service, not systemd unit service)")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// check `unitFileDir`
	// TODO: if invalid dir path, print warning
	err = unix.MkdirIfNotExist(unitFileDir)
	if err != nil {
		return Systemd{}, err
	}

	if !strings.HasSuffix(unitFileDir, "/") {
		// add trailing slash
		unitFileDir += "/"
	}

	service = Systemd{s, unitFileDir}
	return service, nil
}

// Implements iSystemdService
type Systemd struct {
	systemctl   Systemctl
	unitFileDir string
}
