package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) DaemonReload() error {
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", "daemon-reload")
	if err != nil {
		return err
	}

	return nil
}
