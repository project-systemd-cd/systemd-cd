package systemctl

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (s systemctl) Enable(service string, startNow bool) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "service", Value: service}, {Name: "startNow", Value: startNow}}))

	command := []string{"enable"}
	if startNow {
		command = append(command, "--now")
	}
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", command...)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
