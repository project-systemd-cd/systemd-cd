package systemd

import (
	"os"
	"reflect"
	"strings"
	"systemd-cd/domain/logger"
)

// NewService implements iSystemdService
func (s Systemd) NewService(name string, uf UnitFileService, env map[string]string) (us unitService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate systemd unit service")
	logger.Logger().Debugf("< name = %v", name)
	logger.Logger().Tracef("< unitFile = %+v", uf)
	logger.Logger().Tracef("< env = %+v", env)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Tracef("> unitService.Name = %+v", us.Name)
			logger.Logger().Tracef("> unitService.unitFile = %+v", us.unitFile)
			logger.Logger().Tracef("> unitService.Path = %+v", us.Path)
			logger.Logger().Tracef("> unitService.EnvironmentFileValues = %+v", us.EnvironmentFileValues)
			logger.Logger().Debug("END   - Instantiate systemd unit service")
		} else {
			logger.Logger().Error("FAILED - Instantiate systemd unit service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// load unit file
	path := strings.Join([]string{s.unitFileDir, name, ".service"}, "")
	loaded, isGeneratedBySystemdCd, err := s.loadUnitFileSerivce(path)
	if err != nil && !os.IsNotExist(err) {
		// fail
		return unitService{}, err
	}

	if os.IsNotExist(err) {
		// unit file not exists
		// generate `.service` file to `path`
		logger.Logger().Debugf("Create unit file (%s)", path)
		err = s.writeUnitFileService(uf, path)
	} else if isGeneratedBySystemdCd {
		// unit file already exists and file generated by systemd-cd
		if !loaded.Equals(uf) {
			// file has changes
			// update `.service` file to `path`
			logger.Logger().Debugf("Update unit file (%s)", path)
			err = s.writeUnitFileService(uf, path)
		}
	} else {
		// unit file already exists and file not generated by systemd-cd
		err = ErrUnitFileNotManaged
	}
	if err != nil {
		// fail
		return unitService{}, err
	}

	if uf.Service.EnvironmentFile != nil {
		// load env file
		envPath := *uf.Service.EnvironmentFile
		var loaded map[string]string
		var isGeneratedBySystemdCd bool
		loaded, isGeneratedBySystemdCd, err = s.loadEnvFile(envPath)
		if err != nil && !os.IsNotExist(err) {
			// fail
			return unitService{}, err
		}

		if os.IsNotExist(err) {
			// unit file not exists
			// generate env file to `envPath`
			logger.Logger().Debugf("Create env file (%s)", envPath)
			err = s.writeEnvFile(env, envPath)
		} else if isGeneratedBySystemdCd {
			// unit file already exists and file generated by systemd-cd
			if !reflect.DeepEqual(env, loaded) {
				// file has changes
				// update env file to `envPath`
				logger.Logger().Debugf("Update env file (%s)", envPath)
				err = s.writeEnvFile(env, envPath)
			}
		} else {
			// unit file already exists and file not generated by systemd-cd
			err = ErrUnitEnvFileNotManaged
		}
		if err != nil {
			// fail
			return unitService{}, err
		}
	}

	// daemon-reload
	err = s.systemctl.DaemonReload()
	if err != nil {
		return unitService{}, err
	}

	us = unitService{s.systemctl, name, uf, path, env}
	return us, nil
}
