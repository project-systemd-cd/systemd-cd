package systemctl

import (
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/unix"
)

func (s systemctl) Stop(service string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "service", Value: service}}))

	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "stop", service)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
