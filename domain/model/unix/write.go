package unix

import (
	"os"
	"systemd-cd/domain/model/logger"
)

func WriteFile(path string, b []byte) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	// Open file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write
	_, err = f.Write(b)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
