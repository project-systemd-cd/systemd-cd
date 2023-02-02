package unix

import (
	"bytes"
	"os"
)

func ReadFile(path string, b *bytes.Buffer) error {
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
