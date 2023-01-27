package unix

import (
	"bytes"
	"os"
	"systemd-cd/domain/model/logger"
)

func ReadFile(path string, b *bytes.Buffer) error {
	logger.Logger().Tracef("Called:\n\tpath: %v\n\tbuffer", path, b)

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
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
		b.Write(bytes.Trim(data, "\x00"))
	}

	logger.Logger().Trace("Finished")
	return nil
}
