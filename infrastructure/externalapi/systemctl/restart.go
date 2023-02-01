package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) Restart(service string) error {
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "restart", service)
	if err != nil {
		return err
	}

	return nil
}
