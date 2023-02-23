package systemd

import (
	"os"
	"systemd-cd/domain/logger"
)

// DeleteService implements iSystemdService
func (s Systemd) DeleteService(u unitService) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Delete systemd service")
	logger.Logger().Debugf("< unitService.Name = %v", u.Name)
	logger.Logger().Tracef("< unitService = %+v", u.Name)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Delete systemd service")
		} else {
			logger.Logger().Error("FAILED - Delete systemd service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = u.Disable(true)
	if err != nil {
		return err
	}

	// Delete `.service` file
	err = os.Remove(u.Path)

	return err
}
