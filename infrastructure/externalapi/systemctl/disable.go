package systemctl

import (
	"systemd-cd/domain/unix"
)

func (s systemctl) Disable(service string, stopNow bool) error {
	command := []string{"disable"}
	if stopNow {
		command = append(command, "--now")
	}
	command = append(command, service)
	_, _, _, err := unix.Execute(unix.ExecuteOption{}, "systemctl", command...)
	if err != nil {
		return err
	}

	return nil
}
