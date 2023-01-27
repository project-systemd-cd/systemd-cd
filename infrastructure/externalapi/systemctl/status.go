package systemctl

import (
	"strings"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/systemd"
	"systemd-cd/domain/model/unix"
)

func (s systemctl) Status(service string) (systemd.Status, error) {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	exitCode, stdout, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "is-active", service)
	if exitCode != 0 && exitCode != 3 && err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return "", err
	}
	outs := strings.Split(stdout.String(), "\n")
	if len(outs) < 1 {
		logger.Logger().Errorf("Error:\n\terr: %v", systemd.ErrUnitStatusCannotUnmarshal)
		return "", systemd.ErrUnitStatusCannotUnmarshal
	}

	switch outs[0] {
	case "active":
		logger.Logger().Tracef("Finished:\n\tStatus: %v", systemd.StatusRunning)
		return systemd.StatusRunning, nil
	case "inactive":
		logger.Logger().Tracef("Finished:\n\tStatus: %v", systemd.StatusStopped)
		return systemd.StatusStopped, nil
	case "failed":
		logger.Logger().Tracef("Finished:\n\tStatus: %v", systemd.StatusFailed)
		return systemd.StatusFailed, nil
	}

	logger.Logger().Errorf("Error:\n\terr: %v", systemd.ErrUnitStatusCannotUnmarshal)
	return "", systemd.ErrUnitStatusCannotUnmarshal
}
