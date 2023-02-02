package unix

import (
	"os"
)

func WriteFile(path string, b []byte) error {
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
