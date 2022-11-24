package toml

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type child struct {
	Name string
}
type parent struct {
	Child      child
	PointerAry *[]child
}

func (p *parent) Equals(p2 parent) bool {
	if p.Child.Name != p2.Child.Name {
		return false
	}
	if p.PointerAry != nil && p2.PointerAry != nil {
		if len(*p.PointerAry) != len(*p2.PointerAry) {
			return false
		}
		for i := range *p.PointerAry {
			if (*p.PointerAry)[i].Name != (*p2.PointerAry)[i].Name {
				return false
			}
		}
		return true
	} else {
		return p.PointerAry == nil && p2.PointerAry == nil
	}
}

func TestDecodeStruct(t *testing.T) {
	type args struct {
		r io.Reader
		i interface{}
	}

	r := bytes.Buffer{}
	r.Write([]byte(`[[PointerAry]]
	Name = "a"
[Child]
  Name = "test"
`))
	aryTmp := []child{{Name: "a"}}
	var ary *[]child = &aryTmp
	w := parent{Child: child{Name: "test"}, PointerAry: ary}

	r2 := bytes.Buffer{}
	r2.Write([]byte(`PointerAry = [ ]
[Child]
  Name = "test"
`))
	aryTmp2 := []child{}
	var ary2 *[]child = &aryTmp2
	w2 := parent{Child: child{Name: "test"}, PointerAry: ary2}

	r3 := bytes.Buffer{}
	r3.Write([]byte(`[Child]
  Name = "test"
`))
	w3 := parent{Child: child{Name: "test"}, PointerAry: nil}

	tests := []struct {
		name    string
		args    args
		wantW   parent
		wantErr bool
	}{
		{
			name:    "",
			args:    args{r: &r, i: &parent{}},
			wantW:   w,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				r: &r2,
				i: &parent{},
			},
			wantW:   w2,
			wantErr: false,
		}, {
			name: "",
			args: args{
				r: &r3,
				i: &parent{},
			},
			wantW:   w3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.r, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			pGot := tt.args.i.(*parent)
			if !pGot.Equals(tt.wantW) {
				t.Errorf("Decode() = %v, want %v", tt.args.i, tt.wantW)
			}
		})
	}
}

func TestDecodeMap(t *testing.T) {
	type args struct {
		r io.Reader
		i interface{}
	}

	r := bytes.Buffer{}
	r.Write([]byte(`sample = [ ]`))
	var p map[string][]string
	w := map[string][]string{"sample": {}}

	tests := []struct {
		name    string
		args    args
		wantW   map[string][]string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{r: &r, i: &p},
			wantW:   w,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.r, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			pGot := tt.args.i.(*map[string][]string)
			if !reflect.DeepEqual(pGot, &tt.wantW) {
				t.Errorf("Decode() = %v, want %v", pGot, &tt.wantW)
			}
		})
	}
}
