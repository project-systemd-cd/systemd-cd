package toml

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		i interface{}
		o EncodeOption
	}

	type child struct {
		Name string
	}
	type parent struct {
		Child child
	}

	indent := ""
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				i: parent{Child: child{Name: "test"}},
				o: EncodeOption{},
			},
			wantW: `[Child]
  Name = "test"
`,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				i: parent{Child: child{Name: "test"}},
				o: EncodeOption{Indent: &indent},
			},
			wantW: `[Child]
Name = "test"
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Encode(w, tt.args.i, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
