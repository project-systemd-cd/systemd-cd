package unix

import (
	"os"
	"systemd-cd/domain/logger"
)

func MkdirIfNotExist(paths ...string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "paths", Value: paths}}))

	for _, path := range paths {
		_, err := os.ReadDir(path)
		if err != nil {
			if os.IsNotExist(err) {
				// if dir not exists, mkdir
				err = os.MkdirAll(path, 0644)
				if err != nil {
					logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
					return err
				}
				logger.Logger().Debugf("Dir %s created.")
			} else {
				// unhandled errors
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
