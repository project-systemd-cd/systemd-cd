package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

type EncodeOption struct {
	Indent *string
}

func Encode(w io.Writer, i interface{}, o EncodeOption) error {
	e := toml.NewEncoder(w)
	if o.Indent != nil {
		e.Indent = *o.Indent
	}
	return e.Encode(i)
}
