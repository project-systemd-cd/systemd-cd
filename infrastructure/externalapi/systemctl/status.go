package systemctl

import (
	"strings"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/systemd"
	"systemd-cd/domain/model/unix"
)

func (s systemctl) Status(service string) (systemd.Status, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "service", Value: service}}))

	exitCode, stdout, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "is-active", service)
	if exitCode != 0 && exitCode != 3 && err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	outs := strings.Split(stdout.String(), "\n")
	if len(outs) < 1 {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: systemd.ErrUnitStatusCannotUnmarshal}}))
		return "", systemd.ErrUnitStatusCannotUnmarshal
	}

	switch outs[0] {
	case "active":
		logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "Status", Value: systemd.StatusRunning}}))
		return systemd.StatusRunning, nil
	case "inactive":
		logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "Status", Value: systemd.StatusStopped}}))
		return systemd.StatusStopped, nil
	case "failed":
		logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "Status", Value: systemd.StatusFailed}}))
		return systemd.StatusFailed, nil
	}

	logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: systemd.ErrUnitStatusCannotUnmarshal}}))
	return "", systemd.ErrUnitStatusCannotUnmarshal
}
