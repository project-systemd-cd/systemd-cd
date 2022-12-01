package systemctl

import (
	"errors"
	"strings"
	"systemd-cd/domain/model/systemd"
)

func New() systemd.Systemctl {
	return systemctl{}
}

type systemctl struct{}

func (s systemctl) DaemonReload() error {
	_, _, stderr, err := executeCommand("systemctl", "daemon-reload")
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Enable(service string, startNow bool) error {
	command := []string{"enable"}
	if startNow {
		command = append(command, "--now")
	}
	_, _, stderr, err := executeCommand("systemctl", command...)
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Disable(service string, stopNow bool) error {
	command := []string{"disable"}
	if stopNow {
		command = append(command, "--now")
	}
	_, _, stderr, err := executeCommand("systemctl", command...)
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Start(service string) error {
	_, _, stderr, err := executeCommand("systemctl", "start", service)
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Stop(service string) error {
	_, _, stderr, err := executeCommand("systemctl", "stop", service)
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Restart(service string) error {
	_, _, stderr, err := executeCommand("systemctl", "restart", service)
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (s systemctl) Status(service string) (systemd.Status, error) {
	_, stdout, stderr, err := executeCommand("systemctl", "is-active", service)
	if err != nil {
		return "", errors.New(stderr.String())
	}
	outs := strings.Split(stdout.String(), "\n")
	if len(outs) < 1 {
		return "", systemd.ErrUnitStatusCannotUnmarshal
	}
	if outs[0] == "active" {
		return systemd.StatusRunning, nil
	}
	if outs[0] == "inactive" {
		return systemd.StatusStopped, nil
	}
	if outs[0] == "failed" {
		return systemd.StatusFailed, nil
	}
	return "", systemd.ErrUnitStatusCannotUnmarshal
}
