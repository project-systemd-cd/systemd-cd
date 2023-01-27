package unix

import (
	"os"
	"systemd-cd/domain/model/logger"
)

func WriteFile(path string, b []byte) error {
	logger.Logger().Tracef("Called:\n\tpath: %v", path)

	// Open file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write
	_, err = f.Write(b)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
