package unix

import (
	"os"
)

func MkdirIfNotExist(paths ...string) error {
	for _, path := range paths {
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
	}

	return nil
}
