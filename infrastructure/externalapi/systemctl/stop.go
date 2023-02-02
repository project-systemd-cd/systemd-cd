package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) Stop(service string) error {
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "stop", service)
	if err != nil {
		return err
	}

	return nil
}
