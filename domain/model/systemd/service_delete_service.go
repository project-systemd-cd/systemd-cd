package systemd

import (
	"os"
	"systemd-cd/domain/model/logger"
)

// DeleteService implements iSystemdService
func (s Systemd) DeleteService(u UnitService) error {
	logger.Logger().Tracef("Called:\n\tu: %v", u)
	err := u.Disable(true)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Delete `.service` file
	err = os.Remove(u.Path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
