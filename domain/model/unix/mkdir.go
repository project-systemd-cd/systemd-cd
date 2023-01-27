package unix

import (
	"os"
	"systemd-cd/domain/model/logger"
)

func MkdirIfNotExist(path string) error {
	logger.Logger().Tracef("Called:\n\tpath: %s", path)

	_, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if dir not exists, mkdir
			err = os.MkdirAll(path, 0644)
			if err != nil {
				logger.Logger().Errorf("Error:\n\terr: %v", err)
				return err
			}
			logger.Logger().Debugf("Dir %s created.")
		} else {
			// unhandled errors
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
