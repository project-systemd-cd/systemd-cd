package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) Start(service string) error {
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "start", service)
	if err != nil {
		return err
	}

	return nil
}
