package systemd

import (
	"bytes"
	"os"
)

func mkdirIfNotExist(path string) error {
	_, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if dir not exists, mkdir
			err = os.MkdirAll(path, 0644)
			if err != nil {
				return err
			}
		} else {
			// unhandled errors
			return err
		}
	}

	return nil
}

func readFile(path string, b *bytes.Buffer) error {
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
			return err
		}
		b.Write(bytes.Trim(data, "\x00"))
	}

	return nil
}

func writeFile(path string, b []byte) error {
	// Open file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
