package systemctl

import "systemd-cd/domain/model/logger"

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
