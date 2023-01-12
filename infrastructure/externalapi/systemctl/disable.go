package systemctl

import "systemd-cd/domain/model/logger"

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
