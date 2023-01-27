package unix

import (
	"bytes"
	"os"
	"systemd-cd/domain/logger"
)

func ReadFile(path string, b *bytes.Buffer) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}, {Value: b}}))

	// Open file
	f, err := os.Open(path)
	if err != nil {
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

	logger.Logger().Trace("Finished")
	return nil
}
