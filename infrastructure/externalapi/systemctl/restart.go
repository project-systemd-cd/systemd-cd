package systemctl

import (
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/unix"
)

func (s systemctl) Restart(service string) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "restart", service)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}