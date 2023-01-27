package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

func Decode(r io.Reader, i interface{}) error {
	_, err := toml.NewDecoder(r).Decode(i)
	return err
}
