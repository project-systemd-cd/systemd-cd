package systemctl

import (
	"strings"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/systemd"
)

func New() systemd.Systemctl {
	return systemctl{}
}

type systemctl struct{}

func (s systemctl) DaemonReload() error {
	logger.Logger().Trace("Called")

	_, _, _, err := executeCommand("systemctl", "daemon-reload")
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v")
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Enable(service string, startNow bool) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v\n\targ.startNow: %v", service, startNow)

	command := []string{"enable"}
	if startNow {
		command = append(command, "--now")
	}
	_, _, _, err := executeCommand("systemctl", command...)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v")
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Disable(service string, stopNow bool) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v\n\targ.stopNow: %v", service, stopNow)

	command := []string{"disable"}
	if stopNow {
		command = append(command, "--now")
	}
	_, _, _, err := executeCommand("systemctl", command...)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Start(service string) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	_, _, _, err := executeCommand("systemctl", "start", service)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Stop(service string) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	_, _, _, err := executeCommand("systemctl", "stop", service)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Restart(service string) error {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	_, _, _, err := executeCommand("systemctl", "restart", service)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}

func (s systemctl) Status(service string) (systemd.Status, error) {
	logger.Logger().Tracef("Called:\n\targ.service: %v", service)

	_, stdout, _, err := executeCommand("systemctl", "is-active", service)
	if err != nil {
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
