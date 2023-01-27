package systemd

import (
	"os"
	"systemd-cd/domain/logger"
)

// DeleteService implements iSystemdService
func (s Systemd) DeleteService(u UnitService) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: u}}))
	err := u.Disable(true)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Delete `.service` file
	err = os.Remove(u.Path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
