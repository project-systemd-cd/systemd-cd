package errors

import (
	"errors"
	"testing"
)

func TestErrNotFound_As(t *testing.T) {
	type fields struct {
		Object string
		Id     string
	}
	type args struct {
		t interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				Object: "User",
				Id:     "1",
			},
			args: args{
				t: ErrNotFound{},
			},
			want: true,
		},
		{
			name: "true",
			fields: fields{
				Object: "User",
				Id:     "1",
			},
			args: args{
				t: &ErrNotFound{},
			},
			want: true,
		},
		{
			name: "true",
			fields: fields{
				Object: "User",
				Id:     "1",
			},
			args: args{
				t: func() any { e := &ErrNotFound{}; return &e }(),
			},
			want: true,
		},
		{
			name: "false",
			fields: fields{
				Object: "User",
				Id:     "1",
			},
			args: args{
				t: &ErrValidation{},
			},
			want: false,
		},
		{
			name: "false",
			fields: fields{
				Object: "User",
				Id:     "1",
			},
			args: args{
				t: func() any { e := errors.New(""); return &e },
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrNotFound{
				Object: tt.fields.Object,
				Id:     tt.fields.Id,
			}
			if got := e.As(tt.args.t); got != tt.want {
				t.Errorf("ErrNotFound.As() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrNotFound_Error(t *testing.T) {
	type fields struct {
		Object string
		Id     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "not found 1",
			fields: fields{
				Object: "Name",
				Id:     "1",
			},
			want: "\"Name\" not found (id: \"1\")",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrNotFound{
				Object: tt.fields.Object,
				Id:     tt.fields.Id,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrNotFound.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
