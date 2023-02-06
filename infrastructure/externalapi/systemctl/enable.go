package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) Enable(service string, startNow bool) error {
	command := []string{"enable"}
	if startNow {
		command = append(command, "--now")
	}
	command = append(command, service)
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", command...)
	if err != nil {
		return err
	}

	return nil
}
