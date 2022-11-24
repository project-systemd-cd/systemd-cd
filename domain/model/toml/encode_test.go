package toml

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		i interface{}
	}

	type child struct {
		Name string
	}
	type parent struct {
		Child child
	}

	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "",
			args: args{i: parent{Child: child{Name: "test"}}},
			wantW: `[Child]
  Name = "test"
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Encode(w, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
