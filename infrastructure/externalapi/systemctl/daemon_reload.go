package systemctl

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (s systemctl) DaemonReload() error {
	logger.Logger().Trace("Called")

	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "daemon-reload")
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
