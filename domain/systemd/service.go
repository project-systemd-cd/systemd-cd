package systemd

import (
	"errors"
	"strings"
	"systemd-cd/domain/logger"
)

var (
	ErrUnitFileNotManaged    = errors.New("unit file not managed by systemd-cd")
	ErrUnitEnvFileNotManaged = errors.New("unit env file not managed by systemd-cd")
)

type IService interface {
	// Generate unit-file.
	// If unit-file already exists, replace it.
	NewService(name string, uf UnitFileService, env map[string]string) (UnitService, error)
	DeleteService(u UnitService) error

	loadUnitFileSerivce(path string) (u UnitFileService, isGeneratedBySystemdCd bool, err error)
	writeUnitFileService(u UnitFileService, path string) error

	loadEnvFile(path string) (e map[string]string, isGeneratedBySystemdCd bool, err error)
	writeEnvFile(e map[string]string, path string) error
}

func New(s Systemctl, unitFileDir string) (IService, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: s}, {Name: "unitFileDir", Value: unitFileDir}}))

	// check `unitFileDir`
	// TODO: if invalid dir path, print warning
	err := mkdirIfNotExist(unitFileDir)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return Systemd{}, err
	}

	if !strings.HasSuffix(unitFileDir, "/") {
		// add trailing slash
		unitFileDir += "/"
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: Systemd{s, unitFileDir}}}))
	return Systemd{s, unitFileDir}, nil
}

// Implements iSystemdService
type Systemd struct {
	systemctl   Systemctl
	unitFileDir string
}
