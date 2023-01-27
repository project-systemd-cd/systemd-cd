package systemctl

import (
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/unix"
)

func (s systemctl) DaemonReload() error {
	logger.Logger().Trace("Called")

	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "daemon-reload")
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v")
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
