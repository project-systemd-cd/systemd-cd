package systemd

import (
	"bytes"
	"os"
	"systemd-cd/domain/model/logger"
)

func mkdirIfNotExist(path string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	_, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if dir not exists, mkdir
			err = os.MkdirAll(path, 0644)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		} else {
			// unhandled errors
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}

func readFile(path string, b *bytes.Buffer) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	// Open file
	f, err := os.Open(path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
		b.Write(bytes.Trim(data, "\x00"))
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: b.String()}}))
	return nil
}

func writeFile(path string, b []byte) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}, {Name: "b", Value: string(b)}}))

	// Open file
	f, err := os.Create(path)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
