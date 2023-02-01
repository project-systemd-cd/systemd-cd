package systemd

import (
	"os"
	"systemd-cd/domain/logger"
)

// DeleteService implements iSystemdService
func (s Systemd) DeleteService(u UnitService) (err error) {
	logger.Logger().Debug("START - Delete systemd service")
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Delete systemd service")
		} else {
			logger.Logger().Error("FAILED - Delete systemd service")
			logger.Logger().Error(err)
		}
	}()

	err = u.Disable(true)
	if err != nil {
		return err
	}

	// Delete `.service` file
	err = os.Remove(u.Path)

	return err
}
