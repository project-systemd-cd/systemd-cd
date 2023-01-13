package systemd

import (
	"bytes"
	"os"
	"systemd-cd/domain/model/logger"
)

func mkdirIfNotExist(path string) error {
	logger.Logger().Tracef("Called:\n\targ.path: %v", path)

	_, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if dir not exists, mkdir
			err = os.MkdirAll(path, 0644)
			if err != nil {
				logger.Logger().Errorf("Error:\n\terr: %v", err)
				return err
			}
		} else {
			// unhandled errors
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}

func readFile(path string, b *bytes.Buffer) error {
	logger.Logger().Tracef("Called:\n\targ.path: %v", path)

	// Open file
	f, err := os.Open(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}
	defer f.Close()

	// Read
	for {
		data := make([]byte, 64)
		var i int
		i, err = f.Read(data)
		if i == 0 {
			break
		}
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
		b.Write(bytes.Trim(data, "\x00"))
	}

	logger.Logger().Tracef("Finished:\n\tstring: %v", b.String())
	return nil
}

func writeFile(path string, b []byte) error {
	logger.Logger().Tracef("Called:\n\targ.path: %v\n\targ.b: %v", path, string(b))

	// Open file
	f, err := os.Create(path)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
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
