package systemctl

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (s systemctl) Restart(service string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "service", Value: service}}))

	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "restart", service)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
