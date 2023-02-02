package systemctl

import (
	"strings"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

func (s systemctl) Status(service string) (systemd.Status, error) {
	exitCode, stdout, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "is-active", service)
	if exitCode != 0 && exitCode != 3 && err != nil {
		return "", err
	}
	outs := strings.Split(stdout.String(), "\n")
	if len(outs) < 1 {
		return "", systemd.ErrUnitStatusCannotUnmarshal
	}

	switch outs[0] {
	case "active":
		return systemd.StatusRunning, nil
	case "inactive":
		return systemd.StatusStopped, nil
	case "failed":
		return systemd.StatusFailed, nil
	}

	return "", systemd.ErrUnitStatusCannotUnmarshal
}
