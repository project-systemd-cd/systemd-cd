package systemctl

import "systemd-cd/domain/model/logger"

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
