package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

func Encode(w io.Writer, i interface{}) error {
	return toml.NewEncoder(w).Encode(i)
}
